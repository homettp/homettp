package main

import (
	"github.com/homettp/homettp/internal/cmd"
	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/support-go/cli"
)

func main() {
	(&cli.App{
		Name:    "Homettp",
		Version: "1.0.0",
		Groups: []*cli.Group{
			&cli.Group{
				Name:  "make",
				Usage: "Make commands",
				Commands: []*cli.Command{
					&cli.Command{
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
			&cli.Group{
				Name:  "web",
				Usage: "Web commands",
				Commands: []*cli.Command{
					&cli.Command{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: cmd.WebServe,
					},
				},
			},
		},
	}).Execute()
}