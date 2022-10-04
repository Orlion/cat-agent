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

func NewMessageTree() *MessageTree {
	return &MessageTree{
		discard: true,
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

func (tree *MessageTree) SetMessage(message Message) {
	tree.message = message
}

func (tree *MessageTree) SetMessageId(messageId string) {
	tree.messageId = messageId
}

func (tree *MessageTree) SetParentMessageId(parentMessageId string) {
	tree.parentMessageId = parentMessageId
}

func (tree *MessageTree) SetRootMessageId(rootMessageId string) {
	tree.rootMessageId = rootMessageId
}

func (tree *MessageTree) SetThreadGroupName(threadGroupName string) {
	tree.threadGroupName = threadGroupName
}

func (tree *MessageTree) SetThreadId(threadId string) {
	tree.threadId = threadId
}

func (tree *MessageTree) SetThreadName(threadName string) {
	tree.threadName = threadName
}

func (tree *MessageTree) SetDiscard(discard bool) {
	tree.discard = discard
}
