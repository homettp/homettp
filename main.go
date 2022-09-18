package main

import (
	"github.com/homettp/homettp/internal/cmd"
	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/support-go/cli"
)

func main() {
	(&cli.App{
		Name:    "Homettp",
		Version: "master",
		Groups: []*cli.Group{
			{
				Name:  "make",
				Usage: "Make commands",
				Commands: []*cli.Command{
					{
						Name:  "user",
						Usage: "Make a user",
						Arguments: []string{
							"username",
							"email",
							"password",
						},
						HandleFunc: cmd.MakeUser,
					},
				},
			},
			{
				Name:  "web",
				Usage: "Web commands",
				Commands: []*cli.Command{
					{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: cmd.WebServe,
					},
				},
			},
		},
	}).Execute()
}
