package models

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

const redisTransactionAttempts = 10

var updateHashScript = redis.NewScript(1, `
if redis.call("EXISTS", KEYS[1]) == 0 then
	return 0
end

redis.call("HSET", KEYS[1], unpack(ARGV))

return 1
`)

// redisTransaction type.
type redisTransaction struct {
	conn     redis.Conn
	commands []redisCommand
}

// redisCommand type.
type redisCommand struct {
	name string
	args []any
}

// newRedisTransaction function.
func newRedisTransaction(conn redis.Conn) *redisTransaction {
	return &redisTransaction{
		conn: conn,
	}
}

// Queue function.
func (rt *redisTransaction) Queue(command string, args ...any) {
	rt.commands = append(rt.commands, redisCommand{name: command, args: args})
}

// CheckUnique function.
func (rt *redisTransaction) CheckUnique(key, value string, id int, errDuplicate error) error {
	score, err := redis.Int(rt.conn.Do("ZSCORE", key, value))
	if errors.Is(err, redis.ErrNil) {
		return nil
	}

	if err != nil {
		return err
	}

	if score != id {
		return errDuplicate
	}

	return nil
}

// Run function.
func (rt *redisTransaction) Run(keys []string, prepare func() error) error {
	for range redisTransactionAttempts {
		rt.commands = nil

		_, err := rt.conn.Do("WATCH", redis.Args{}.AddFlat(keys)...)
		if err != nil {
			return err
		}

		err = prepare()
		if err != nil {
			_, _ = rt.conn.Do("UNWATCH")

			return err
		}

		err = rt.conn.Send("MULTI")
		if err != nil {
			return err
		}

		for _, command := range rt.commands {
			err = rt.conn.Send(command.name, command.args...)
			if err != nil {
				return err
			}
		}

		replies, err := redis.Values(rt.conn.Do("EXEC"))
		if errors.Is(err, redis.ErrNil) {
			continue
		}

		if err != nil {
			return err
		}

		for _, reply := range replies {
			if replyErr, ok := reply.(redis.Error); ok {
				return replyErr
			}
		}

		return nil
	}

	return ErrConflict
}
