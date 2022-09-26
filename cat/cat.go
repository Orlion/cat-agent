package cat

import (
	"fmt"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
)

var cat *Cat

type Cat struct {
	manager      *Manager
	msgIdFactory *MessageIdFactory
}

func (cat *Cat) run() {
	cat.manager.run()
}

func (cat *Cat) shutdown() {
	config.Shutdown()
	cat.manager.shutdown()
}

func Init(conf *config.Config) error {
	if err := config.Init(conf); err != nil {
		return err
	}

	cat = &Cat{
		manager:      newManager(),
		msgIdFactory: newMessageIdFactory(),
	}

	return nil
}

func Flush(tree *message.MessageTree) {
	cat.manager.flush(tree)
}

func GetNextId(domain string) (string, error) {
	configDomain := config.GetInstance().GetDomain()
	if domain != configDomain {
		return "", fmt.Errorf("cat-agent's domain is %s, not %s", configDomain, domain)
	}

	return getNextId(), nil
}

func getNextId() string {
	return cat.msgIdFactory.getNextId()
}

func Shutdown() {
	cat.shutdown()
}
