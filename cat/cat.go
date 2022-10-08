package cat

import (
	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
)

var catInstance *Cat

type Cat struct {
	inShutdown   bool
	manager      *Manager
	msgIdFactory *MessageIdFactory
}

func (cat *Cat) run() {
	cat.msgIdFactory.run()
	cat.manager.run()
}

func (cat *Cat) shutdown() {
	log.Info("cat shutdown...")
	cat.inShutdown = true
	config.Shutdown()
	cat.manager.shutdown()
	log.Info("cat exit")
}

func (cat *Cat) shuttingDown() bool {
	return cat.inShutdown
}

func (cat *Cat) send(tree *message.MessageTree) {
	if cat.shuttingDown() || !config.GetInstance().IsEnabled() {
		return
	}

	cat.manager.send(tree)
}

func (cat *Cat) createMessageId(domain string) []byte {
	return cat.msgIdFactory.getNextId(domain)
}

func Init(conf *config.Config) error {
	if err := config.Init(conf); err != nil {
		return err
	}

	catInstance = &Cat{
		manager:      newManager(),
		msgIdFactory: newMessageIdFactory(),
	}

	catInstance.run()

	return nil
}

func Send(tree *message.MessageTree) {
	catInstance.send(tree)
}

func CreateMessageId(domain string) []byte {
	return catInstance.createMessageId(domain)
}

func Shutdown() {
	catInstance.shutdown()
}
