package message

type MessageTree struct {
	Domain          string
	Hostname        string
	IpAddress       string
	message         Message
	MessageId       string
	ParentMessageId string
	RootMessageId   string
	threadGroupName string
	threadId        string
	threadName      string
	discard         bool
	hitSample       bool
}

func NewMessageTree() *MessageTree {
	return &MessageTree{}
}

func (tree *MessageTree) GetMessage() Message {
	return tree.message
}

func (tree *MessageTree) canDiscard() bool {
	return true
}
