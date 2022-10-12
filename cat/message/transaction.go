package message

type Transaction struct {
	baseMessage
	children         []Message
	durationInMicros int64
}

func NewTransaction(t, name, status, data string, timestampInMillis int64, children []Message, durationInMicros int64) *Transaction {
	return &Transaction{
		baseMessage:      newBaseMessage(t, name, status, data, timestampInMillis),
		children:         children,
		durationInMicros: durationInMicros,
	}
}

func (trans *Transaction) GetChildren() []Message {
	return trans.children
}

func (trans *Transaction) GetDurationInMicros() int64 {
	return trans.durationInMicros
}

func (trans *Transaction) SetDurationInMicros(durationInMicros int64) {
	trans.durationInMicros = durationInMicros
}

func (trans *Transaction) AddChild(child Message) {
	trans.children = append(trans.children, child)
}
