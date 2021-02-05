package models

type CallStatus string

const (
	Pending    CallStatus = "pending"
	InProgress CallStatus = "in_progress"
	Completed  CallStatus = "completed"
	Failed     CallStatus = "failed"
)

type Call struct {
	Id        int64      `json:"id" redis:"id"`
	CommandId int        `json:"command_id" redis:"command_id"`
	Status    CallStatus `json:"status" redis:"status"`
	Payload   string     `json:"payload" redis:"payload"`
	Output    string     `json:"output" redis:"output"`
	CreatedAt Timestamp  `json:"created_at" redis:"created_at"`
}

type CallRepository interface {
	Create(*Call) error
	Find(int64) (*Call, error)
	FindAllByCommand(*Command) ([]*Call, error)
	Update(*Call, *Call) error
	Delete(*Call) error
}

func NewCall(command *Command) *Call {
	return &Call{
		CommandId: command.Id,
		Status:    Pending,
	}
}
