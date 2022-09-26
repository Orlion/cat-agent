package cat

import "github.com/Orlion/cat-agent/cat/message"

type LocalAggregator struct {
	ta *TransactionAggregator
	ea *EventAggregator
}

func newLocalAggregator() *LocalAggregator {
	return &LocalAggregator{
		ta: newTransactionAggregator(),
		ea: newEventAggregator(),
	}
}

func (la *LocalAggregator) aggregate(tree *message.MessageTree) {
	msg := tree.GetMessage()
	switch msg.(type) {
	case *message.Transaction:
		la.analyzerProcessTransaction(msg.(*message.Transaction))
	case *message.Event:
	default:
	}
}

func (la *LocalAggregator) analyzerProcessTransaction(transaction *message.Transaction) {
	la.ta.logTransaction(transaction)

	for _, child := range transaction.GetChildren() {
		switch child.(type) {
		case *message.Transaction:
			la.analyzerProcessTransaction(child.(*message.Transaction))
		case *message.Event:
			la.ea.logEvent(child.(*message.Event))
		}
	}
}

func computeDuration(durationInMillis int) int {
	if durationInMillis < 1 {
		return 1
	} else if durationInMillis < 20 {
		return durationInMillis
	} else if durationInMillis < 200 {
		return durationInMillis - durationInMillis%5
	} else if durationInMillis < 500 {
		return durationInMillis - durationInMillis%20
	} else if durationInMillis < 2000 {
		return durationInMillis - durationInMillis%50
	} else if durationInMillis < 20000 {
		return durationInMillis - durationInMillis%500
	} else if durationInMillis < 1000000 {
		return durationInMillis - durationInMillis%10000
	} else {
		dk := 524288
		if durationInMillis > 3600*1000 {
			dk = 3600 * 1000
		} else {
			for dk < durationInMillis {
				dk <<= 1
			}
		}
		return dk
	}
}
