package message

type Heartbeat struct {
	baseMessage
}

func NewHeartbeat(t, name, status, data string, timestampInMillis int64) *Heartbeat {
	return &Heartbeat{
		baseMessage: newBaseMessage(t, name, status, data, timestampInMillis),
	}
}
