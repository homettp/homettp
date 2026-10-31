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

	if command.StartsWithPayload() {
		return ErrInvalidValue
	}

	id, err := redis.Int(conn.Do("INCR", rcr.RedisKeyPrefix+commandKeyPrefix+commandIDKey))
	if err != nil {
		return err
	}

	created := *command
	created.ID = id
	created.Token = token
	created.CreatedAt = Timestamp(time.Now())

	nameKey := rcr.RedisKeyPrefix + commandKeyPrefix + commandNameKey
	name := strings.ToLower(created.Name)

	tx := newRedisTransaction(conn)

	err = tx.Run([]string{nameKey}, func() error {
		err := tx.CheckUnique(nameKey, name, id, ErrDuplicateName)
		if err != nil {
			return err
		}

		tx.Queue("HSET", redis.Args{}.Add(rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(id)).AddFlat(&created)...)
		tx.Queue("ZADD", nameKey, id, name)

		return nil
	})
	if err != nil {
		return err
	}

	*command = created

	return nil
}

// Find function.
func (rcr *RedisCommandRepository) Find(id int) (*Command, error) {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	return rcr.find(conn, id)
}

// find function.
func (rcr *RedisCommandRepository) find(conn redis.Conn, id int) (*Command, error) {
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

	commands := make([]*Command, 0, len(names))
	var id int

	for _, name := range names {
		id, err = strconv.Atoi(values[name])
		if err != nil {
			return nil, err
		}

		command, err := rcr.find(conn, id)
		if errors.Is(err, ErrNoRecord) {
			continue
		}

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

	if newCommand.StartsWithPayload() {
		return ErrInvalidValue
	}

	key := rcr.RedisKeyPrefix + commandKeyPrefix + strconv.Itoa(command.ID)
	nameKey := rcr.RedisKeyPrefix + commandKeyPrefix + commandNameKey
	newName := strings.ToLower(newCommand.Name)

	var updated *Command

	tx := newRedisTransaction(conn)

	err := tx.Run([]string{key, nameKey}, func() error {
		current, err := rcr.find(conn, command.ID)
		if err != nil {
			return err
		}

		oldName := strings.ToLower(current.Name)

		if oldName != newName {
			err = tx.CheckUnique(nameKey, newName, current.ID, ErrDuplicateName)
			if err != nil {
				return err
			}

			tx.Queue("ZREM", nameKey, oldName)
			tx.Queue("ZADD", nameKey, current.ID, newName)
		}

		current.Name = newCommand.Name
		current.Image = newCommand.Image
		current.Value = newCommand.Value

		tx.Queue("HSET", redis.Args{}.Add(key).AddFlat(current)...)
		updated = current

		return nil
	})
	if err != nil {
		return err
	}

	*command = *updated

	return nil
}

// UpdateToken function.
func (rcr *RedisCommandRepository) UpdateToken(command *Command, token string) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	updated, err := redis.Bool(updateHashScript.Do(
		conn, rcr.RedisKeyPrefix+commandKeyPrefix+strconv.Itoa(command.ID), "token", token,
	))
	if err != nil {
		return err
	}

	if !updated {
		return ErrNoRecord
	}

	command.Token = token

	return nil
}

// Delete function.
func (rcr *RedisCommandRepository) Delete(command *Command) error {
	conn := rcr.RedisPool.Get()
	defer conn.Close()

	key := rcr.RedisKeyPrefix + commandKeyPrefix + strconv.Itoa(command.ID)
	callsKey := rcr.RedisKeyPrefix + commandKeyPrefix + callKeyPrefix + strconv.Itoa(command.ID)
	tx := newRedisTransaction(conn)

	err := tx.Run([]string{key, callsKey}, func() error {
		current, err := rcr.find(conn, command.ID)
		if err != nil {
			return err
		}

		ids, err := redis.Int64s(conn.Do("LRANGE", callsKey, 0, -1))
		if err != nil {
			return err
		}

		for _, id := range ids {
			tx.Queue("DEL", rcr.RedisKeyPrefix+callKeyPrefix+strconv.FormatInt(id, 10))
		}

		tx.Queue("DEL", callsKey)
		tx.Queue("ZREM", rcr.RedisKeyPrefix+commandKeyPrefix+commandNameKey, strings.ToLower(current.Name))
		tx.Queue("DEL", key)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
