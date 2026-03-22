package config

import (
	"os"
	"strconv"

	"github.com/petaki/support-go/cli"
)

// Config type.
type Config struct {
	Debug               bool
	Addr                string
	URL                 string
	Key                 string
	RedisURL            string
	RedisKeyPrefix      string
	CommandTimeout      int
	CommandWorkerCount  int
	CommandHistoryLimit int
}

// NewConfig function.
func NewConfig(command *cli.Command, arguments []string) (*Config, error) {
	debug := command.FlagSet().Bool("debug", false, "Application Debug Mode")
	addr := command.FlagSet().String("addr", os.Getenv("APP_ADDR"), "Application Address")
	url := command.FlagSet().String("url", os.Getenv("APP_URL"), "Application URL")
	key := command.FlagSet().String("key", os.Getenv("APP_KEY"), "Application Key")
	redisURL := command.FlagSet().String("redis-url", os.Getenv("REDIS_URL"), "Redis URL")
	redisKeyPrefix := command.FlagSet().String("redis-key-prefix", os.Getenv("REDIS_KEY_PREFIX"), "Redis Key Prefix")

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

	_, err = command.Parse(arguments)
	if err != nil {
		return nil, err
	}

	return &Config{
		Debug:               *debug,
		Addr:                *addr,
		URL:                 *url,
		Key:                 *key,
		RedisURL:            *redisURL,
		RedisKeyPrefix:      *redisKeyPrefix,
		CommandTimeout:      *commandTimeout,
		CommandWorkerCount:  *commandWorkerCount,
		CommandHistoryLimit: *commandHistoryLimit,
	}, nil
}
