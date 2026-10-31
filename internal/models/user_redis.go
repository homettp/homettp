package models

import (
	"crypto/subtle"
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

	password, err := bcrypt.GenerateFromPassword(user.Password, 12)
	if err != nil {
		return err
	}

	id, err := redis.Int(conn.Do("INCR", rur.RedisKeyPrefix+userKeyPrefix+userIDKey))
	if err != nil {
		return err
	}

	created := *user
	created.ID = id
	created.Password = password
	created.CreatedAt = Timestamp(time.Now())

	usernameKey := rur.RedisKeyPrefix + userKeyPrefix + userUsernameKey
	emailKey := rur.RedisKeyPrefix + userKeyPrefix + userEmailKey
	username := strings.ToLower(created.Username)
	email := strings.ToLower(created.Email)

	tx := newRedisTransaction(conn)

	err = tx.Run([]string{usernameKey, emailKey}, func() error {
		err := tx.CheckUnique(usernameKey, username, id, ErrDuplicateUsername)
		if err != nil {
			return err
		}

		err = tx.CheckUnique(emailKey, email, id, ErrDuplicateEmail)
		if err != nil {
			return err
		}

		tx.Queue("HSET", redis.Args{}.Add(rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(id)).AddFlat(&created)...)
		tx.Queue("ZADD", usernameKey, id, username)
		tx.Queue("ZADD", emailKey, id, email)

		return nil
	})
	if err != nil {
		return err
	}

	*user = created

	return nil
}

// Find function.
func (rur *RedisUserRepository) Find(id int) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	return rur.find(conn, id)
}

// find function.
func (rur *RedisUserRepository) find(conn redis.Conn, id int) (*User, error) {
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

	users := make([]*User, 0, len(usernames))

	for _, username := range usernames {
		id, err := strconv.Atoi(values[username])
		if err != nil {
			return nil, err
		}

		user, err := rur.find(conn, id)
		if errors.Is(err, ErrNoRecord) {
			continue
		}

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

	var password []byte
	var err error

	if len(newUser.Password) > 0 {
		password, err = bcrypt.GenerateFromPassword(newUser.Password, 12)
		if err != nil {
			return err
		}
	}

	key := rur.RedisKeyPrefix + userKeyPrefix + strconv.Itoa(user.ID)
	usernameKey := rur.RedisKeyPrefix + userKeyPrefix + userUsernameKey
	emailKey := rur.RedisKeyPrefix + userKeyPrefix + userEmailKey
	newUsername := strings.ToLower(newUser.Username)
	newEmail := strings.ToLower(newUser.Email)

	var updated *User

	tx := newRedisTransaction(conn)

	err = tx.Run([]string{key, usernameKey, emailKey}, func() error {
		current, err := rur.find(conn, user.ID)
		if err != nil {
			return err
		}

		oldUsername := strings.ToLower(current.Username)
		oldEmail := strings.ToLower(current.Email)

		if oldUsername != newUsername {
			err = tx.CheckUnique(usernameKey, newUsername, current.ID, ErrDuplicateUsername)
			if err != nil {
				return err
			}

			tx.Queue("ZREM", usernameKey, oldUsername)
			tx.Queue("ZADD", usernameKey, current.ID, newUsername)
		}

		if oldEmail != newEmail {
			err = tx.CheckUnique(emailKey, newEmail, current.ID, ErrDuplicateEmail)
			if err != nil {
				return err
			}

			tx.Queue("ZREM", emailKey, oldEmail)
			tx.Queue("ZADD", emailKey, current.ID, newEmail)
		}

		current.Username = newUser.Username
		current.Email = newUser.Email
		current.IsEnabled = newUser.IsEnabled

		if password != nil {
			current.Password = password
			current.RememberToken = ""
		}

		tx.Queue("HSET", redis.Args{}.Add(key).AddFlat(current)...)
		updated = current

		return nil
	})
	if err != nil {
		return err
	}

	*user = *updated

	return nil
}

// UpdateRememberToken function.
func (rur *RedisUserRepository) UpdateRememberToken(user *User, token string) error {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	updated, err := redis.Bool(updateHashScript.Do(
		conn, rur.RedisKeyPrefix+userKeyPrefix+strconv.Itoa(user.ID), "remember_token", token,
	))
	if err != nil {
		return err
	}

	if !updated {
		return ErrNoRecord
	}

	user.RememberToken = token

	return nil
}

// Authenticate function.
func (rur *RedisUserRepository) Authenticate(usernameOrEmail, password string) (*User, error) {
	conn := rur.RedisPool.Get()
	defer conn.Close()

	if len(password) > maxPasswordBytes {
		return nil, rejectCredentials(password)
	}

	usernameOrEmail = strings.ToLower(usernameOrEmail)

	id, err := redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, usernameOrEmail))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			id, err = redis.Int(conn.Do("ZSCORE", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, usernameOrEmail))
			if errors.Is(err, redis.ErrNil) {
				return nil, rejectCredentials(password)
			}

			if err != nil {
				return nil, err
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
		return nil, rejectCredentials(password)
	}

	var user User

	err = redis.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	if !user.IsEnabled {
		return nil, rejectCredentials(password)
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

	if user.RememberToken == "" || subtle.ConstantTimeCompare([]byte(user.RememberToken), []byte(token)) != 1 {
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

	key := rur.RedisKeyPrefix + userKeyPrefix + strconv.Itoa(user.ID)
	tx := newRedisTransaction(conn)

	err := tx.Run([]string{key}, func() error {
		current, err := rur.find(conn, user.ID)
		if err != nil {
			return err
		}

		tx.Queue("ZREM", rur.RedisKeyPrefix+userKeyPrefix+userUsernameKey, strings.ToLower(current.Username))
		tx.Queue("ZREM", rur.RedisKeyPrefix+userKeyPrefix+userEmailKey, strings.ToLower(current.Email))
		tx.Queue("DEL", key)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
