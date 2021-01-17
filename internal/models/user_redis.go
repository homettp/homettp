package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

const (
	userKeyPrefix   = "user:"
	userIdKey       = "id"
	userUsernameKey = "username"
	userEmailKey    = "email"
)

type RedisUserRepository struct {
	RedisPool      *redis.Pool
	RedisKeyPrefix string
}

func (rur *RedisUserRepository) Create(user *User) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, user.Username))
	if err == nil {
		return ErrDuplicateUsername
	} else if !errors.Is(err, redis.ErrNil) {
		return err
	}

	_, err = redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, user.Email))
	if err == nil {
		return ErrDuplicateEmail
	} else if !errors.Is(err, redis.ErrNil) {
		return err
	}

	user.Password, err = bcrypt.GenerateFromPassword(user.Password, 12)
	if err != nil {
		return err
	}

	user.Id, err = redis.Int(conn.Do("INCR", rur.RedisKeyPrefix+userKeyPrefix+userIdKey))
	if err != nil {
		return err
	}

	user.CreatedAt = Timestamp(time.Now())

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send(
		"HMSET", redis.Args{}.Add(rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.Id)).AddFlat(user)...,
	)
	if err != nil {
		return err
	}

	err = conn.Send("ZADD", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, "NX", user.Id, user.Username)
	if err != nil {
		return err
	}

	err = conn.Send("ZADD", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, "NX", user.Id, user.Email)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

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

func (rur *RedisUserRepository) UpdateRememberToken(user *User, token string) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.Id), "remember_token", token)
	if err != nil {
		return err
	}

	user.RememberToken = token

	return nil
}

func (rur *RedisUserRepository) Authenticate(usernameOrEmail, password string) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	id, err := redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, usernameOrEmail))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			id, err = redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, usernameOrEmail))

			if errors.Is(err, redis.ErrNil) {
				return nil, ErrInvalidCredentials
			}
		}

		return nil, err
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
		} else {
			return nil, err
		}
	}

	return &user, nil
}

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
