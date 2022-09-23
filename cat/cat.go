package cat

import (
	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/manager"
	"github.com/Orlion/cat-agent/cat/message"
)

var cat *Cat

type Cat struct {
	manager *manager.Manager
}

func Init(conf *config.Config) error {
	if err := config.Init(conf); err != nil {
		return err
	}
	cat = &Cat{
		manager: manager.NewManager(),
	}

	return nil
}

func Flush(tree *message.MessageTree) {
	cat.manager.Flush(tree)
}
