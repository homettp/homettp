package cmd

import (
	"fmt"
	"strconv"

	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
	"github.com/homettp/homettp/internal/service"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/forms"
)

// MakeUser command.
func MakeUser(group *cli.Group, command *cli.Command, arguments []string) int {
	appConfig, err := config.NewConfig(command, arguments)
	if err != nil {
		cli.ErrorLog.Fatal(err)

		return command.PrintHelp(group)
	}

	parsed := command.FlagSet().Args()

	form := forms.New(map[string]any{
		"username":   parsed[0],
		"email":      parsed[1],
		"password":   parsed[2],
		"is_enabled": true,
	})

	models.UserCreateRules(form)

	if !form.IsValid() {
		return command.PrintError(fmt.Errorf("make user: invalid arguments: %v", form.Errors))
	}

	redisPool := service.RedisPool(appConfig)
	defer redisPool.Close()

	userRepository := &models.RedisUserRepository{
		RedisPool:      redisPool,
		RedisKeyPrefix: appConfig.RedisKeyPrefix,
	}

	user := (&models.User{}).Fill(form)
	err = userRepository.Create(user)

	if err != nil {
		return command.PrintError(err)
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
