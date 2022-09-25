package cat

import (
	"fmt"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/manager"
	"github.com/Orlion/cat-agent/cat/message"
)

var cat *Cat

type Cat struct {
	manager      *manager.Manager
	msgIdFactory *MessageIdFactory
}

func (cat *Cat) run() {
}

func Init(conf *config.Config) error {
	if err := config.Init(conf); err != nil {
		return err
	}

	cat = &Cat{
		manager:      manager.NewManager(),
		msgIdFactory: newMessageIdFactory(),
	}

	return nil
}

func Flush(tree *message.MessageTree) {
	cat.manager.Flush(tree)
}

func GetNextId(domain string) (string, error) {
	configDomain := config.GetInstance().GetDomain()
	if domain != configDomain {
		return "", fmt.Errorf("cat-agent's domain is %s, not %s", configDomain, domain)
	}

	return cat.msgIdFactory.getNextId(domain), nil
}
