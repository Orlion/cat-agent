package message

type Transaction struct {
	baseMessage
	children            []Message
	rawDurationInMicros int64
}

func NewTransaction(t, name, status, data string, timestampInMillis int64, children []Message, rawDurationInMicros int64) *Transaction {
	return &Transaction{
		baseMessage:         newBaseMessage(t, name, status, data, timestampInMillis),
		children:            children,
		rawDurationInMicros: rawDurationInMicros,
	}
}

func (trans *Transaction) GetChildren() []Message {
	return trans.children
}

func (trans *Transaction) GetRawDurationInMicros() int64 {
	return trans.rawDurationInMicros
}

func (trans *Transaction) AddChild(child Message) {
	trans.children = append(trans.children, child)
}
