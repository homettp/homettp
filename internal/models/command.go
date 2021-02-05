package models

import (
	"fmt"

	"github.com/petaki/support-go/forms"
)

type CommandImage string

const (
	Door            CommandImage = "door"
	Light           CommandImage = "light"
	Outlet          CommandImage = "outlet"
	Plug            CommandImage = "plug"
	Sensor          CommandImage = "sensor"
	PayloadVariable string       = "{payload}"
)

type Command struct {
	Id        int          `json:"id" redis:"id"`
	Name      string       `json:"name" redis:"name"`
	Token     string       `json:"-" redis:"token"`
	Image     CommandImage `json:"image" redis:"image"`
	Value     string       `json:"value" redis:"value"`
	CreatedAt Timestamp    `json:"created_at" redis:"created_at"`
}

type CommandRepository interface {
	Create(*Command, string) error
	Find(int) (*Command, error)
	FindAll() ([]*Command, error)
	Update(*Command, *Command) error
	UpdateToken(*Command, string) error
	Delete(*Command) error
}

func NewCommand() *Command {
	return &Command{
		Image: Light,
	}
}

func CommandCreateRules(form *forms.Form) {
	form.Required("name", "image", "value")
}

func CommandUpdateRules(form *forms.Form) {
	form.Required("name", "image", "value")
}

func (c *Command) Fill(form *forms.Form) *Command {
	c.Name = form.Data["name"].(string)
	c.Image = CommandImage(form.Data["image"].(string))
	c.Value = form.Data["value"].(string)

	return c
}

func (c *Command) Path(url string) string {
	return fmt.Sprintf("%s/call?id=%v&token=%s", url, c.Id, c.Token)
}
