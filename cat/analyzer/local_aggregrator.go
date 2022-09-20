package analyzer

import "github.com/Orlion/cat-agent/cat/message"

type LocalAggregator struct {
	transaction *TransactionAggregator
}

func (aggregator *LocalAggregator) Aggregate(tree *message.MessageTree) {
	msg := tree.GetMessage()
	switch msg.(type) {
	case *message.Transaction:
	case *message.Event:
	default:
	}
}
