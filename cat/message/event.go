package message

type Event struct {
	baseMessage
}

func NewEvent() *Event {
	return &Event{}
}
