package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/forms"
)

func makeUser(group *cli.Group, command *cli.Command, arguments []string) int {
	enabled := command.FlagSet().Bool("enabled", true, "User Enabled")
	redisURL, redisKeyPrefix := createRedisFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	form := forms.New(map[string]interface{}{
		"username":   parsed[0],
		"email":      parsed[1],
		"password":   parsed[2],
		"is_enabled": *enabled,
	})

	models.UserCreateRules(form)

	if !form.IsValid() {
		return printError(fmt.Errorf("make user: invalid arguments: %v", form.Errors))
	}

	userRepository := &models.RedisUserRepository{
		RedisPool:      newRedisPool(*redisURL),
		RedisKeyPrefix: *redisKeyPrefix,
	}

	user := (&models.User{}).Fill(form)
	err = userRepository.Create(user)

	if err != nil {
		return printError(err)
	}

	return (&cli.Table{
		Headers: []string{"ID", "Username", "Email", "IsEnabled"},
		Rows: [][]string{
			{
				strconv.Itoa(user.ID),
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
