package models

type CallStatus string

const (
	Pending    CallStatus = "pending"
	InProgress CallStatus = "in_progress"
	Completed  CallStatus = "completed"
	Failed     CallStatus = "failed"
)

type Call struct {
	Id        int        `json:"id" redis:"id"`
	CommandId int        `json:"command_id" redis:"command_id"`
	Status    CallStatus `json:"status" redis:"status"`
	Output    string     `json:"output" redis:"output"`
	CreatedAt Timestamp  `json:"created_at" redis:"created_at"`
}
