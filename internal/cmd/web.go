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
		cli.ErrorLog.Fatal(err)

		return command.PrintHelp(group)
	}

	redisPool := service.RedisPool(appConfig)
	defer redisPool.Close()

	web.Serve(appConfig, redisPool)

	return cli.Success
}
