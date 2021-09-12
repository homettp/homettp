package cmd

import (
	"os"

	"github.com/homettp/homettp/internal/web"
	"github.com/petaki/support-go/cli"
)

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	key := command.FlagSet().String("key", os.Getenv("APP_KEY"), "Application Key")

	redisURL, redisKeyPrefix := createRedisFlags(command)
	commandTimeout, commandWorkerCount, commandHistoryLimit := createCommandFlags(command)

	redisPool := newRedisPool(*redisURL)
	defer redisPool.Close()

	web.Serve(
		*debug,
		*addr,
		*url,
		*key,
		*redisKeyPrefix,
		redisPool,
		*commandTimeout,
		*commandWorkerCount,
		*commandHistoryLimit,
	)

	return cli.Success
}
