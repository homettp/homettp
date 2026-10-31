package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	callKeyPrefix = "call:"
	callIDKey     = "id"
)

// RedisCallRepository type.
type RedisCallRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
	HistoryLimit   int
}

// Create function.
func (rcr *RedisCallRepository) Create(call *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	var err error

	call.ID, err = redis.Int64(conn.Do("INCR", rcr.RedisKeyPrefix+callKeyPrefix+callIDKey))
	if err != nil {
		return err
	}

	call.CreatedAt = Timestamp(time.Now())

	commandKey := rcr.RedisKeyPrefix + commandKeyPrefix + strconv.Itoa(call.CommandID)

	tx := newRedisTransaction(conn)

	err = tx.Run([]string{commandKey}, func() error {
		exists, err := redis.Bool(conn.Do("EXISTS", commandKey))
		if err != nil {
			return err
		}

		if !exists {
			return ErrNoRecord
		}

		tx.Queue("HSET", redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.ID, 10)).AddFlat(call)...)
		tx.Queue("LPUSH", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), call.ID)

		return nil
	})
	if err != nil {
		return err
	}

	ids, err := redis.Int64s(
		conn.Do("LRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), rcr.HistoryLimit, -1),
	)
	if err != nil {
		return err
	}

	var finishedIDs []int64

	for _, id := range ids {
		status, err := redis.String(conn.Do("HGET", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10), "status"))
		if err != nil && !errors.Is(err, redis.ErrNil) {
			return err
		}

		if !CallStatus(status).IsRunning() {
			finishedIDs = append(finishedIDs, id)
		}
	}

	if len(finishedIDs) == 0 {
		return nil
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, id := range finishedIDs {
		err = conn.Send("DEL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}

		err = conn.Send(
			"LREM", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), 0, id,
		)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// Find function.
func (rcr *RedisCallRepository) Find(id int64) (*Call, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	return rcr.find(conn, id)
}

// find function.
func (rcr *RedisCallRepository) find(conn redis.Conn, id int64) (*Call, error) {
	values, err := redis.Values(
		conn.Do("HGETALL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10)),
	)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrNoRecord
	}

	var call Call

	err = redis.ScanStruct(values, &call)
	if err != nil {
		return nil, err
	}

	return &call, nil
}

// FindAllByCommand function.
func (rcr *RedisCallRepository) FindAllByCommand(command *Command) ([]*Call, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	ids, err := redis.Int64s(
		conn.Do("LRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(command.ID), 0, -1),
	)
	if err != nil {
		return nil, err
	}

	calls := make([]*Call, 0, len(ids))

	for _, id := range ids {
		call, err := rcr.find(conn, id)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				continue
			}

			return nil, err
		}

		calls = append(calls, call)
	}

	return calls, nil
}

// Update function.
func (rcr *RedisCallRepository) Update(call, newCall *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	call.Status = newCall.Status
	call.Output = newCall.Output

	updated, err := redis.Bool(updateHashScript.Do(
		conn, redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.ID, 10)).AddFlat(call)...,
	))
	if err != nil {
		return err
	}

	if !updated {
		return ErrNoRecord
	}

	return nil
}

// Delete function.
func (rcr *RedisCallRepository) Delete(call *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"LREM", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), 0, call.ID,
	)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.ID, 10))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
