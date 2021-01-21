package models

type CommandImage string

const (
	Door   CommandImage = "door"
	Light  CommandImage = "light"
	Outlet CommandImage = "outlet"
	Plug   CommandImage = "plug"
	Rgb    CommandImage = "rgb"
	Window CommandImage = "window"
)

type Command struct {
	Id        int          `json:"id" redis:"id"`
	Name      string       `json:"name" redis:"name"`
	Token     string       `json:"token" redis:"token"`
	Image     CommandImage `json:"image" redis:"image"`
	Timeout   int          `json:"timeout" redis:"timeout"`
	Value     string       `json:"value" redis:"value"`
	CreatedAt Timestamp    `json:"created_at" redis:"created_at"`
}
