package cat

import (
	"fmt"

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

func (cat *Cat) flush(tree *message.MessageTree) {
	if cat.shuttingDown() {
		return
	}

	cat.manager.flush(tree)
}

func (cat *Cat) getNextId() string {
	return cat.msgIdFactory.getNextId()
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

func Flush(tree *message.MessageTree) {
	catInstance.flush(tree)
}

func GetNextId(domain string) (string, error) {
	configDomain := config.GetInstance().GetDomain()
	if domain != configDomain {
		return "", fmt.Errorf("cat-agent's domain is %s, not %s", configDomain, domain)
	}

	return getNextId(), nil
}

func getNextId() string {
	return catInstance.getNextId()
}

func Shutdown() {
	catInstance.shutdown()
}
