package message

type Event struct {
	baseMessage
}

func NewEvent(t, name, status, data string, timestampInMillis int64) *Event {
	return &Event{
		baseMessage: newBaseMessage(t, name, status, data, timestampInMillis),
	}
}
