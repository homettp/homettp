package models

import "github.com/petaki/support-go/forms"

type CommandImage string

const (
	Door   CommandImage = "door"
	Light  CommandImage = "light"
	Outlet CommandImage = "outlet"
	Plug   CommandImage = "plug"
	Sensor CommandImage = "sensor"
)

type Command struct {
	Id        int          `json:"id" redis:"id"`
	Name      string       `json:"name" redis:"name"`
	Token     string       `json:"-" redis:"token"`
	Image     CommandImage `json:"image" redis:"image"`
	Timeout   int          `json:"timeout" redis:"timeout"`
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
		Image:   Light,
		Timeout: 60,
	}
}

func CommandCreateRules(form *forms.Form) {
	form.Required("name", "image", "timeout", "value")
	form.Min("timeout", 1)
}

func CommandUpdateRules(form *forms.Form) {
	form.Required("name", "image", "timeout", "value")
	form.Min("timeout", 1)
}

func (u *Command) Fill(form *forms.Form) *Command {
	u.Name = form.Data["name"].(string)
	u.Image = CommandImage(form.Data["image"].(string))
	u.Timeout = int(form.Data["timeout"].(float64))
	u.Value = form.Data["value"].(string)

	return u
}
