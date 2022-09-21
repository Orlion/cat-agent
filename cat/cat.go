package cat

import "github.com/Orlion/cat-agent/cat/message"

var cat *Cat

type Cat struct {
	manager *Manager
}

func Init() {
	cat = &Cat{
		manager: newManager(),
	}
}

func Flush(domain string, hostname string, ipAddress string, msg message.Message, messageId string, parentMessageId string, rootMessageId string, threadGroupName string, threadId string, threadName string) {
	cat.manager.flush(domain, hostname, ipAddress, msg, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName)
}

func NewTransaction() *message.Transaction {
	return message.NewTransaction()
}

func NewEvent() *message.Event {
	return message.NewEvent()
}
