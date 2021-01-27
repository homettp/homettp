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

func (rcr *RedisCommandRepository) Find(id int) (*Command, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(id)))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrNoRecord
	}

	var command Command

	err = redis.ScanStruct(values, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func (rcr *RedisCommandRepository) FindAll() ([]*Command, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	names, err := redis.StringMap(conn.Do("ZRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, 0, -1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	var commands []*Command
	var id int

	for _, value := range names {
		id, err = strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		command, err := rcr.Find(id)
		if err != nil {
			return nil, err
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func (rcr *RedisCommandRepository) Update(command, newCommand *Command) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	var err error

	if command.Name != newCommand.Name {
		_, err = redis.Int(
			conn.Do("ZSCORE", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(newCommand.Name)),
		)
		if err == nil {
			return ErrDuplicateName
		} else if !errors.Is(err, redis.ErrNil) {
			return err
		}
	}

	command.Image = newCommand.Image
	command.Timeout = newCommand.Timeout
	command.Value = newCommand.Value

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	if command.Name != newCommand.Name {
		err = conn.Send(
			"ZREM", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(command.Name),
		)
		if err != nil {
			return err
		}

		command.Name = newCommand.Name

		err = conn.Send(
			"ZADD", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, "NX", command.Id, strings.ToLower(command.Name),
		)
		if err != nil {
			return err
		}
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.Id)).AddFlat(command)...,
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

func (rcr *RedisCommandRepository) UpdateToken(command *Command, token string) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.Id), "token", token)
	if err != nil {
		return err
	}

	command.Token = token

	return nil
}

func (rcr *RedisCommandRepository) Delete(command *Command) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZREM", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(command.Name),
	)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.Id))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
