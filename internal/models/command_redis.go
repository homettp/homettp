package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	commandKeyPrefix = "command:"
	commandIdKey     = "id"
	commandNameKey   = "name"
)

type RedisCommandRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

func (rcr *RedisCommandRepository) Create(command *Command) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	_, err := redis.Int(
		conn.Do("ZSCORE", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(command.Name)),
	)
	if err == nil {
		return ErrDuplicateName
	} else if !errors.Is(err, redis.ErrNil) {
		return err
	}

	command.Id, err = redis.Int(conn.Do("INCR", rcr.RedisKeyPrefix+commandKeyPrefix+commandIdKey))
	if err != nil {
		return err
	}

	command.CreatedAt = Timestamp(time.Now())

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.Id)).AddFlat(command)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZADD", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, "NX", command.Id, strings.ToLower(command.Name),
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
