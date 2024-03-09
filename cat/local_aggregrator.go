package cat

import (
	"context"
	"sync"

	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/pkg/atomicx"
)

type LocalAggregator struct {
	ta         *TransactionAggregator
	ea         *EventAggregator
	cancel     func()
	wg         *sync.WaitGroup
	inShutdown atomicx.Bool
}

func newLocalAggregator() *LocalAggregator {
	return &LocalAggregator{
		ta: newTransactionAggregator(),
		ea: newEventAggregator(),
		wg: new(sync.WaitGroup),
	}
}

func (la *LocalAggregator) run() {
	var ctx context.Context
	ctx, la.cancel = context.WithCancel(context.Background())
	la.wg.Add(2)
	go func() {
		la.ea.run(ctx)
		la.wg.Done()
	}()
	go func() {
		la.ta.run(ctx)
		la.wg.Done()
	}()
}

func (la *LocalAggregator) shutdown() {
	la.inShutdown.SetTrue()
	la.cancel()
	la.wg.Wait()
}

func (la *LocalAggregator) aggregate(tree *message.MessageTree) {
	if la.inShutdown.Get() {
		return
	}

	domain := string(tree.GetDomain())

	msg := tree.GetMessage()
	switch msg.(type) {
	case *message.Transaction:
		la.analyzerProcessTransaction(domain, msg.(*message.Transaction))
	case *message.Event:
		la.ea.logEvent(domain, msg.(*message.Event))
	default:
	}
}

func (la *LocalAggregator) analyzerProcessTransaction(domain string, transaction *message.Transaction) {
	la.ta.logTransaction(domain, transaction)
	for _, child := range transaction.GetChildren() {
		switch child.(type) {
		case *message.Transaction:
			la.analyzerProcessTransaction(domain, child.(*message.Transaction))
		case *message.Event:
			la.ea.logEvent(domain, child.(*message.Event))
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
