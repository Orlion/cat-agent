package sender

import (
	"github.com/Orlion/cat-agent/cat/message"
)

type Sender interface {
	Offer(tree *message.MessageTree)
	Run()
	Shutdown()
}
