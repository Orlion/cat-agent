package message

type MessageTree struct {
	domain          string
	hostname        string
	ipAddress       string
	message         Message
	messageId       string
	parentMessageId string
	rootMessageId   string
	threadGroupName string
	threadId        string
	threadName      string
}

func NewMessageTree(domain, hostname, ipAddress string, message Message, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName string) *MessageTree {
	return &MessageTree{
		domain:          domain,
		hostname:        hostname,
		ipAddress:       ipAddress,
		message:         message,
		messageId:       messageId,
		parentMessageId: parentMessageId,
		rootMessageId:   rootMessageId,
		threadGroupName: threadGroupName,
		threadId:        threadId,
		threadName:      threadName,
	}
}

func (tree *MessageTree) GetDomain() string {
	return tree.domain
}

func (tree *MessageTree) GetHostname() string {
	return tree.hostname
}

func (tree *MessageTree) GetIpAddress() string {
	return tree.ipAddress
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
