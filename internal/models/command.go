package models

import (
	"fmt"
	"strings"

	"github.com/petaki/support-go/forms"
)

// CommandImage type.
type CommandImage string

const (
	// Door image.
	Door CommandImage = "door"

	// Light image.
	Light CommandImage = "light"

	// Outlet image.
	Outlet CommandImage = "outlet"

	// Plug image.
	Plug CommandImage = "plug"

	// Sensor image.
	Sensor CommandImage = "sensor"

	// PayloadVariable pattern.
	PayloadVariable string = "%p"
)

// Command type.
type Command struct {
	ID        int          `json:"id" redis:"id"`
	Name      string       `json:"name" redis:"name"`
	Token     string       `json:"-" redis:"token"`
	Image     CommandImage `json:"image" redis:"image"`
	Value     string       `json:"value" redis:"value"`
	CreatedAt Timestamp    `json:"created_at" redis:"created_at"`
}

// CommandRepository type.
type CommandRepository interface {
	Create(*Command, string) error
	Find(int) (*Command, error)
	FindAll() ([]*Command, error)
	Update(*Command, *Command) error
	UpdateToken(*Command, string) error
	Delete(*Command) error
}

// NewCommand function.
func NewCommand() *Command {
	return &Command{
		Image: Light,
	}
}

// IsValid function.
func (ci CommandImage) IsValid() bool {
	switch ci {
	case Door, Light, Outlet, Plug, Sensor:
		return true
	default:
		return false
	}
}

// CommandCreateRules function.
func CommandCreateRules(form *forms.Form) {
	form.Required("name", "image", "value")
	isString(form, "name", "image", "value")
	isCommandImage(form, "image")
}

// CommandUpdateRules function.
func CommandUpdateRules(form *forms.Form) {
	form.Required("name", "image", "value")
	isString(form, "name", "image", "value")
	isCommandImage(form, "image")
}

// Fill function.
func (c *Command) Fill(form *forms.Form) *Command {
	c.Name, _ = form.Data["name"].(string)
	image, _ := form.Data["image"].(string)
	c.Image = CommandImage(image)
	c.Value, _ = form.Data["value"].(string)

	return c
}

// Path function.
func (c *Command) Path(url string) string {
	return fmt.Sprintf("%s/call?commandID=%v&token=%s", url, c.ID, c.Token)
}

// StartsWithPayload function.
func (c *Command) StartsWithPayload() bool {
	return strings.HasPrefix(strings.TrimSpace(c.Value), PayloadVariable)
}
