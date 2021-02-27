package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/petaki/support-go/cli"
)

func createRedisFlags(command *cli.Command) (*string, *string) {
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")
	redisKeyPrefix := command.FlagSet().String("redis-key-prefix", os.Getenv("REDIS_KEY_PREFIX"), "Redis Key Prefix")

	return redisURL, redisKeyPrefix
}

func createCommandFlags(command *cli.Command) (*int, *int, *int) {
	envCommandTimeout, err := strconv.Atoi(os.Getenv("COMMAND_TIMEOUT"))
	if err != nil {
		envCommandTimeout = 60
	}

	envCommandWorkerCount, err := strconv.Atoi(os.Getenv("COMMAND_WORKER_COUNT"))
	if err != nil {
		envCommandWorkerCount = 2
	}

	envCommandHistoryLimit, err := strconv.Atoi(os.Getenv("COMMAND_HISTORY_LIMIT"))
	if err != nil {
		envCommandHistoryLimit = 20
	}

	commandTimeout := command.FlagSet().Int("command-timeout", envCommandTimeout, "Command Timeout")
	commandWorkerCount := command.FlagSet().Int("command-worker-count", envCommandWorkerCount, "Command Worker Count")
	commandHistoryLimit := command.FlagSet().Int("command-history-limit", envCommandHistoryLimit, "Command History Limit")

	return commandTimeout, commandWorkerCount, commandHistoryLimit
}

func printError(err error) int {
	fmt.Println(err)

	return 1
}
