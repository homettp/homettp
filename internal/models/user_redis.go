package models

import (
	"bytes"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

const (
	userKeyPrefix   = "user:"
	userIDKey       = "id"
	userUsernameKey = "username"
	userEmailKey    = "email"
)

// RedisUserRepository type.
type RedisUserRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

// Create function.
func (rur *RedisUserRepository) Create(user *User) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	_, err := redis.Int(
		conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, strings.ToLower(user.Username)),
	)
	if err == nil {
		return ErrDuplicateUsername
	} else if !errors.Is(err, redis.ErrNil) {
		return err
	}

	_, err = redis.Int(
		conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, strings.ToLower(user.Email)),
	)
	if err == nil {
		return ErrDuplicateEmail
	} else if !errors.Is(err, redis.ErrNil) {
		return err
	}

	user.Password, err = bcrypt.GenerateFromPassword(user.Password, 12)
	if err != nil {
		return err
	}

	user.ID, err = redis.Int(conn.Do("INCR", rur.RedisKeyPrefix+userKeyPrefix+userIDKey))
	if err != nil {
		return err
	}

	user.CreatedAt = Timestamp(time.Now())

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.ID)).AddFlat(user)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZADD", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, "NX", user.ID, strings.ToLower(user.Username),
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZADD", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, "NX", user.ID, strings.ToLower(user.Email),
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
func (rur *RedisUserRepository) Find(id int) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(id)))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrNoRecord
	}

	var user User

	err = redis.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindAll function.
func (rur *RedisUserRepository) FindAll() ([]*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	values, err := redis.StringMap(
		conn.Do("ZRANGE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, 0, -1, "WITHSCORES"),
	)
	if err != nil {
		return nil, err
	}

	usernames := make([]string, 0, len(values))

	for key := range values {
		usernames = append(usernames, key)
	}

	sort.Strings(usernames)

	var users []*User
	var id int

	for _, username := range usernames {
		id, err = strconv.Atoi(values[username])
		if err != nil {
			return nil, err
		}

		user, err := rur.Find(id)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Update function.
func (rur *RedisUserRepository) Update(user, newUser *User) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	var err error

	if user.Username != newUser.Username {
		_, err = redis.Int(
			conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, strings.ToLower(newUser.Username)),
		)
		if err == nil {
			return ErrDuplicateUsername
		} else if !errors.Is(err, redis.ErrNil) {
			return err
		}
	}

	if user.Email != newUser.Email {
		_, err = redis.Int(
			conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, strings.ToLower(newUser.Email)),
		)
		if err == nil {
			return ErrDuplicateEmail
		} else if !errors.Is(err, redis.ErrNil) {
			return err
		}
	}

	if len(newUser.Password) > 0 && !bytes.Equal(user.Password, newUser.Password) {
		user.Password, err = bcrypt.GenerateFromPassword(newUser.Password, 12)
		if err != nil {
			return err
		}
	}

	user.IsEnabled = newUser.IsEnabled

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	if user.Username != newUser.Username {
		err = conn.Send(
			"ZREM", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, strings.ToLower(user.Username),
		)
		if err != nil {
			return err
		}

		user.Username = newUser.Username

		err = conn.Send(
			"ZADD", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, "NX", user.ID, strings.ToLower(user.Username),
		)
		if err != nil {
			return err
		}
	}

	if user.Email != newUser.Email {
		err = conn.Send(
			"ZREM", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, strings.ToLower(user.Email),
		)
		if err != nil {
			return err
		}

		user.Email = newUser.Email

		err = conn.Send(
			"ZADD", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, "NX", user.ID, strings.ToLower(user.Email),
		)
		if err != nil {
			return err
		}
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.ID)).AddFlat(user)...,
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

// UpdateRememberToken function.
func (rur *RedisUserRepository) UpdateRememberToken(user *User, token string) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.ID), "remember_token", token)
	if err != nil {
		return err
	}

	user.RememberToken = token

	return nil
}

// Authenticate function.
func (rur *RedisUserRepository) Authenticate(usernameOrEmail, password string) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	usernameOrEmail = strings.ToLower(usernameOrEmail)

	id, err := redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, usernameOrEmail))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			id, err = redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, usernameOrEmail))

			if errors.Is(err, redis.ErrNil) {
				return nil, ErrInvalidCredentials
			}
		} else {
			return nil, err
		}
	}

	values, err := redis.Values(conn.Do("HGETALL", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(id)))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrInvalidCredentials
	}

	var user User

	err = redis.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	if !user.IsEnabled {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	return &user, nil
}

// AuthenticateByRememberToken function.
func (rur *RedisUserRepository) AuthenticateByRememberToken(id int, token string) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(id)))
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, ErrInvalidCredentials
	}

	var user User

	err = redis.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	if user.RememberToken == "" || user.RememberToken != token {
		return nil, ErrInvalidCredentials
	}

	if !user.IsEnabled {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// Delete function.
func (rur *RedisUserRepository) Delete(user *User) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZREM", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, strings.ToLower(user.Username),
	)
	if err != nil {
		return err
	}

	err = conn.Send(
		"ZREM", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, strings.ToLower(user.Email),
	)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.ID))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
