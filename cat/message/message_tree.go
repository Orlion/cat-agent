package message

type MessageTree struct {
	message         Message
	messageId       string
	parentMessageId string
	rootMessageId   string
	threadGroupName string
	threadId        string
	threadName      string
	discard         bool
}

func NewMessageTree(message Message, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName string, discard bool) *MessageTree {
	return &MessageTree{
		message:         message,
		messageId:       messageId,
		parentMessageId: parentMessageId,
		rootMessageId:   rootMessageId,
		threadGroupName: threadGroupName,
		threadId:        threadId,
		threadName:      threadName,
		discard:         discard,
	}
}

func (tree *MessageTree) GetMessage() Message {
	return tree.message
}

func (tree *MessageTree) GetMessageId() string {
	return tree.messageId
}

func (tree *MessageTree) GetParentMessageId() string {
	return tree.parentMessageId
}

func (tree *MessageTree) GetRootMessageId() string {
	return tree.rootMessageId
}

func (tree *MessageTree) GetThreadGroupName() string {
	return tree.threadGroupName
}

func (tree *MessageTree) GetThreadId() string {
	return tree.threadId
}

func (tree *MessageTree) GetThreadName() string {
	return tree.threadName
}

func (tree *MessageTree) CanDiscard() bool {
	return tree.discard
}
