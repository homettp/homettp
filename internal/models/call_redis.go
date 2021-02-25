package models

import (
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
	RedisPool           *redis.Pool
	RedisKeyPrefix      string
	CommandHistoryLimit int
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

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.ID, 10)).AddFlat(call)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"LPUSH", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), call.ID,
	)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	ids, err := redis.Int64s(
		conn.Do("LRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), rcr.CommandHistoryLimit, -1),
	)
	if err != nil {
		return err
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = conn.Send("DEL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}
	}

	err = conn.Send(
		"LTRIM", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandID), 0, rcr.CommandHistoryLimit-1,
	)
	if err != nil {
		return err
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

	var calls []*Call

	for _, id := range ids {
		call, err := rcr.Find(id)
		if err != nil {
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

	_, err := conn.Do(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.ID, 10)).AddFlat(call)...,
	)
	if err != nil {
		return err
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
