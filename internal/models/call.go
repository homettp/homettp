package models

// CallStatus type.
type CallStatus string

const (
	// Pending status.
	Pending CallStatus = "pending"

	// InProgress status.
	InProgress CallStatus = "in_progress"

	// Completed status.
	Completed CallStatus = "completed"

	// Failed status.
	Failed CallStatus = "failed"
)

// Call type.
type Call struct {
	ID        int64      `json:"id" redis:"id"`
	CommandID int        `json:"command_id" redis:"command_id"`
	Status    CallStatus `json:"status" redis:"status"`
	Payload   string     `json:"payload" redis:"payload"`
	Output    string     `json:"output" redis:"output"`
	CreatedAt Timestamp  `json:"created_at" redis:"created_at"`
}

// CallRepository type.
type CallRepository interface {
	Create(*Call) error
	Find(int64) (*Call, error)
	FindAllByCommand(*Command) ([]*Call, error)
	Update(*Call, *Call) error
	Delete(*Call) error
}

// NewCall function.
func NewCall(command *Command) *Call {
	return &Call{
		CommandID: command.ID,
		Status:    Pending,
	}
}
