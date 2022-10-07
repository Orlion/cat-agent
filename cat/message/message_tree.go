package message

type MessageTree struct {
	message         Message
	domain          []byte
	messageId       []byte
	parentMessageId []byte
	rootMessageId   []byte
	threadGroupName []byte
	threadId        []byte
	threadName      []byte
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

func (tree *MessageTree) GetDomain() []byte {
	return tree.domain
}

func (tree *MessageTree) GetMessageId() []byte {
	return tree.messageId
}

func (tree *MessageTree) GetParentMessageId() []byte {
	return tree.parentMessageId
}

func (tree *MessageTree) GetRootMessageId() []byte {
	return tree.rootMessageId
}

func (tree *MessageTree) GetThreadGroupName() []byte {
	return tree.threadGroupName
}

func (tree *MessageTree) GetThreadId() []byte {
	return tree.threadId
}

func (tree *MessageTree) GetThreadName() []byte {
	return tree.threadName
}

func (tree *MessageTree) CanDiscard() bool {
	return tree.discard
}

func (tree *MessageTree) SetMessage(message Message) {
	tree.message = message
}

func (tree *MessageTree) SetDomain(domain []byte) {
	tree.domain = domain
}

func (tree *MessageTree) SetMessageId(messageId []byte) {
	tree.messageId = messageId
}

func (tree *MessageTree) SetParentMessageId(parentMessageId []byte) {
	tree.parentMessageId = parentMessageId
}

func (tree *MessageTree) SetRootMessageId(rootMessageId []byte) {
	tree.rootMessageId = rootMessageId
}

func (tree *MessageTree) SetThreadGroupName(threadGroupName []byte) {
	tree.threadGroupName = threadGroupName
}

func (tree *MessageTree) SetThreadId(threadId []byte) {
	tree.threadId = threadId
}

func (tree *MessageTree) SetThreadName(threadName []byte) {
	tree.threadName = threadName
}

func (tree *MessageTree) SetDiscard(discard bool) {
	tree.discard = discard
}
