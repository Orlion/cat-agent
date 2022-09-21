package message

type Transaction struct {
	baseMessage
	children        []Message
	durationInMicro int64
	durationStart   int64
}

func NewTransaction(t, name, status, data string, timestampInMillis int) *Transaction {
	return &Transaction{
		baseMessage: newBaseMessage(t, name, status, data, timestampInMillis),
	}
}

func (trans *Transaction) GetChildren() []Message {
	return trans.children
}

func (trans *Transaction) GetDurationInMillis() int64 {
	return trans.durationInMicro / 1000
}

func (trans *Transaction) AddChild(child Message) {
	trans.children = append(trans.children, child)
}
