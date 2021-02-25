package models

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	commandKeyPrefix = "command:"
	commandIDKey     = "id"
	commandNameKey   = "name"
)

// RedisCommandRepository type.
type RedisCommandRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

// Create function.
func (rcr *RedisCommandRepository) Create(command *Command, token string) error {
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

	if command.Value == PayloadVariable {
		return ErrInvalidValue
	}

	command.ID, err = redis.Int(conn.Do("INCR", rcr.RedisKeyPrefix+commandKeyPrefix+commandIDKey))
	if err != nil {
		return err
	}

	command.Token = token
	command.CreatedAt = Timestamp(time.Now())

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.ID)).AddFlat(command)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZADD", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, "NX", command.ID, strings.ToLower(command.Name),
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

// FindAll function.
func (rcr *RedisCommandRepository) FindAll() ([]*Command, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	values, err := redis.StringMap(
		conn.Do("ZRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, 0, -1, "WITHSCORES"),
	)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(values))

	for key := range values {
		names = append(names, key)
	}

	sort.Strings(names)

	var commands []*Command
	var id int

	for _, name := range names {
		id, err = strconv.Atoi(values[name])
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

// Update function.
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

	if newCommand.Value == PayloadVariable {
		return ErrInvalidValue
	}

	command.Image = newCommand.Image
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
			"ZADD", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, "NX", command.ID, strings.ToLower(command.Name),
		)
		if err != nil {
			return err
		}
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.ID)).AddFlat(command)...,
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

// UpdateToken function.
func (rcr *RedisCommandRepository) UpdateToken(command *Command, token string) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.ID), "token", token)
	if err != nil {
		return err
	}

	command.Token = token

	return nil
}

// Delete function.
func (rcr *RedisCommandRepository) Delete(command *Command) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	calls, err := redis.Int64s(
		conn.Do("LRANGE", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(command.ID), 0, -1),
	)
	if err != nil {
		return err
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, id := range calls {
		err = conn.Send("DEL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}
	}

	err = conn.Send("DEL", rcr.RedisKeyPrefix+commandKeyPrefix+callKeyPrefix+strconv.Itoa(command.ID))
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZREM", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(command.Name),
	)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.ID))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
