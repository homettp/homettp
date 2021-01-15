package forms

type Error struct {
	Status int
	Msg    string
}

func (e *Error) Error() string {
	return e.Msg
}
