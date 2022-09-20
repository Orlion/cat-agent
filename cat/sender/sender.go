package sender

import "github.com/Orlion/cat-agent/cat/message"

type Sender struct {
	normal chan *message.MessageTree
	high   chan *message.MessageTree
}

func (s *Sender) Offer(tree *message.MessageTree) {
	if tree.Message.IsSuccess() {
		select {
		case s.normal <- tree:
		default:

		}
	} else {
		select {
		case s.high <- tree:
		default:

		}
	}
}
