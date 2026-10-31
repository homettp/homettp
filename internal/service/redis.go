package service

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/homettp/homettp/internal/config"
)

const (
	redisTimeout     = 3 * time.Second
	redisPingTimeout = 500 * time.Millisecond
)

// RedisPool function.
func RedisPool(appConfig *config.Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(
				appConfig.RedisURL,
				redis.DialConnectTimeout(redisTimeout),
				redis.DialReadTimeout(redisTimeout),
				redis.DialWriteTimeout(redisTimeout),
			)
		},
		TestOnBorrow: func(conn redis.Conn, lastUsed time.Time) error {
			if time.Since(lastUsed) < time.Minute {
				return nil
			}

			_, err := redis.DoWithTimeout(conn, redisPingTimeout, "PING")

			return err
		},
	}
}
