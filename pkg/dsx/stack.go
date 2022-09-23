package dsx

import "github.com/Orlion/cat-agent/cat/message"

type TransactionStack struct {
	list []*message.Transaction
}

func NewTransactionStack() *TransactionStack {
	return &TransactionStack{
		list: make([]*message.Transaction, 0),
	}
}

func (s *TransactionStack) Push(trans *message.Transaction) {
	s.list = append(s.list, trans)
}

func (s *TransactionStack) Pop() *message.Transaction {
	trans := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return trans
}

func (s *TransactionStack) IsEmpty() bool {
	return len(s.list) == 0
}

func (s *TransactionStack) Peek() *message.Transaction {
	if s.IsEmpty() {
		return nil
	}
	return s.list[len(s.list)-1]
}
