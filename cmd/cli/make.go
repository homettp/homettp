package main

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/cli"
)

func makeUser(group *cli.Group, command *cli.Command, arguments []string) int {
	enabled := command.FlagSet().Bool("enabled", true, "User Enabled")
	redisUrl, redisKeyPrefix := createRedisFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	userRepository := &models.RedisUserRepository{
		RedisPool:      newRedisPool(*redisUrl),
		RedisKeyPrefix: *redisKeyPrefix,
	}

	user := &models.User{
		Username:  parsed[0],
		Email:     parsed[1],
		Password:  []byte(parsed[2]),
		IsEnabled: *enabled,
	}

	err = userRepository.Create(user)

	if err != nil {
		return printError(err)
	}

	return (&cli.Table{
		Headers: []string{"Id", "Username", "Email", "IsEnabled"},
		Rows: [][]string{
			{
				strconv.Itoa(user.Id),
				user.Username,
				user.Email,
				strconv.FormatBool(user.IsEnabled),
			},
		},
	}).Print()
}

func newRedisPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
	}
}
