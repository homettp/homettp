package models

import (
	"strconv"

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

func (rcr *RedisCallRepository) Find(id int64) (*Call, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10)))
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
