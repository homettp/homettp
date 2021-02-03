package main

import (
	"os"
	"strconv"

	"github.com/homettp/homettp/internal/web"
	"github.com/petaki/support-go/cli"
)

func webServe(group *cli.Group, command *cli.Command, arguments []string) int {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	key := command.FlagSet().String("key", os.Getenv("APP_KEY"), "Application Key")

	redisUrl, redisKeyPrefix := createRedisFlags(command)

	envCommandTimeout, err := strconv.Atoi(os.Getenv("COMMAND_TIMEOUT"))
	if err != nil {
		return printError(err)
	}

	commandTimeout := command.FlagSet().Int("command-timeout", envCommandTimeout, "Command Timeout")

	web.Serve(
		*debug,
		*addr,
		*url,
		*key,
		*redisUrl,
		*redisKeyPrefix,
		*commandTimeout,
	)

	return 0
}
