package config

import (
	"errors"
	"testing"

	"github.com/petaki/support-go/cli"
)

func TestNewConfigRejectsInvalidCommandSettings(t *testing.T) {
	for _, test := range []struct {
		env      string
		value    string
		expected error
	}{
		{"COMMAND_TIMEOUT", "0", ErrInvalidCommandTimeout},
		{"COMMAND_TIMEOUT", "5s", cli.ErrInvalidEnvValue},
		{"COMMAND_TIMEOUT", "10000000000", ErrInvalidCommandTimeout},
		{"COMMAND_WORKER_COUNT", "0", ErrInvalidCommandWorkerCount},
		{"COMMAND_WORKER_COUNT", "abc", cli.ErrInvalidEnvValue},
		{"COMMAND_HISTORY_LIMIT", "0", ErrInvalidCommandHistoryLimit},
		{"COMMAND_HISTORY_LIMIT", "abc", cli.ErrInvalidEnvValue},
	} {
		t.Run(test.env+"="+test.value, func(t *testing.T) {
			t.Setenv(test.env, test.value)

			_, err := NewConfig(&cli.Command{Name: "serve"}, nil)
			if !errors.Is(err, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func TestNewConfigDefaults(t *testing.T) {
	t.Setenv("COMMAND_TIMEOUT", "")
	t.Setenv("COMMAND_WORKER_COUNT", "")
	t.Setenv("COMMAND_HISTORY_LIMIT", "")

	appConfig, err := NewConfig(&cli.Command{Name: "serve"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if appConfig.CommandTimeout != 60 || appConfig.CommandWorkerCount != 2 || appConfig.CommandHistoryLimit != 20 {
		t.Fatalf("unexpected defaults %+v", appConfig)
	}
}

func TestNewConfigSecretsFromEnvAndFlags(t *testing.T) {
	t.Setenv("APP_KEY", "env-key")
	t.Setenv("REDIS_URL", "redis://env:6379/0")

	appConfig, err := NewConfig(&cli.Command{Name: "serve"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if appConfig.Key != "env-key" || appConfig.RedisURL != "redis://env:6379/0" {
		t.Fatalf("expected env values, got %q %q", appConfig.Key, appConfig.RedisURL)
	}

	appConfig, err = NewConfig(&cli.Command{Name: "serve"}, []string{"--key", "flag-key", "--redis-url", "redis://flag:6379/0"})
	if err != nil {
		t.Fatal(err)
	}

	if appConfig.Key != "flag-key" || appConfig.RedisURL != "redis://flag:6379/0" {
		t.Fatalf("expected flag values, got %q %q", appConfig.Key, appConfig.RedisURL)
	}
}
