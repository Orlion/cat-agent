package cat

import (
	"github.com/Orlion/cat-agent/cat/manager"
	"github.com/Orlion/cat-agent/cat/message"
)

var cat *Cat

type Cat struct {
	manager *manager.Manager
}

func Init(config *Config) error {
	if err := withDefaultConf(config); err != nil {
		return err
	}
	cat = &Cat{
		manager: manager.NewManager(),
	}

	return nil
}

func Flush(msg message.Message, messageId string, parentMessageId string, rootMessageId string, threadGroupName string, threadId string, threadName string, discard bool) {
	cat.manager.Flush(msg, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName, discard)
}
