package sender

import "github.com/Orlion/cat-agent/cat/message"

type TcpSender struct {
	normal chan *message.MessageTree
	high   chan *message.MessageTree
}

func NewTcpSender() *TcpSender {
	return &TcpSender{}
}

func (s *TcpSender) Offer(tree *message.MessageTree) {
	if tree.GetMessage().IsSuccess() {
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
