package config

import (
	"fmt"
	"math"
	"time"

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
	debug := command.Bool("debug", "APP_DEBUG", false, "Application Debug Mode")
	addr := command.String("addr", "APP_ADDR", ":4000", "Application Address")
	url := command.String("url", "APP_URL", "", "Application URL")
	key := command.Secret("key", "APP_KEY", "Application Key")
	redisURL := command.Secret("redis-url", "REDIS_URL", "Redis URL")
	redisKeyPrefix := command.String("redis-key-prefix", "REDIS_KEY_PREFIX", "", "Redis Key Prefix")
	commandTimeout := command.Int("command-timeout", "COMMAND_TIMEOUT", 60, "Command Timeout")
	commandWorkerCount := command.Int("command-worker-count", "COMMAND_WORKER_COUNT", 2, "Command Worker Count")
	commandHistoryLimit := command.Int("command-history-limit", "COMMAND_HISTORY_LIMIT", 20, "Command History Limit")

	_, err := command.Parse(arguments)
	if err != nil {
		return nil, err
	}

	if *commandTimeout < 1 || int64(*commandTimeout) > math.MaxInt64/int64(time.Second) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidCommandTimeout, *commandTimeout)
	}

	if *commandWorkerCount < 1 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidCommandWorkerCount, *commandWorkerCount)
	}

	if *commandHistoryLimit < 1 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidCommandHistoryLimit, *commandHistoryLimit)
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
