package models

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	callKeyPrefix = "call:"
	callIdKey     = "id"
)

type RedisCallRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

func (rcr *RedisCallRepository) Create(call *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	var err error

	call.Id, err = redis.Int64(conn.Do("INCR", rcr.RedisKeyPrefix+callKeyPrefix+callIdKey))
	if err != nil {
		return err
	}

	call.CreatedAt = Timestamp(time.Now())

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.Id, 10)).AddFlat(call)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"LPUSH", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandId), call.Id,
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

func (rcr *RedisCallRepository) FindAllByCommand(command *Command) ([]*Call, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	ids, err := redis.Int64s(
		conn.Do("LRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(command.Id), 0, -1),
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

func (rcr *RedisCallRepository) Update(call, newCall *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	call.Status = newCall.Status
	call.Output = newCall.Output

	_, err := conn.Do(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(call.Id, 10)).AddFlat(call)...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (rcr *RedisCallRepository) Delete(call *Call) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"LREM", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(call.CommandId), 0, call.Id,
	)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.FormatInt(call.Id, 10))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
