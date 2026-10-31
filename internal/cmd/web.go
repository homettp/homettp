package cmd

import (
	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/service"
	"github.com/homettp/homettp/internal/web"
	"github.com/petaki/support-go/cli"
)

// WebServe command.
func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	appConfig, err := config.NewConfig(command, arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	redisPool := service.RedisPool(appConfig)
	defer redisPool.Close()

	err = web.Serve(appConfig, redisPool)
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}
