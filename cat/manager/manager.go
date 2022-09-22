package manager

import (
	"sync/atomic"

	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/cat/sender"
)

type Manager struct {
	aggregator  *LocalAggregator
	sender      sender.Sender
	sampleCount uint64
}

func NewManager() *Manager {
	manager := &Manager{
		sender: sender.NewTcpSender(),
	}

	manager.aggregator = newLocalAggregator(manager)

	return manager
}

func (m *Manager) Flush(msg message.Message, messageId string, parentMessageId string, rootMessageId string, threadGroupName string, threadId string, threadName string, discard bool) {
	tree := message.NewMessageTree(msg, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName, discard)
	if tree.CanDiscard() && m.isHitSample() {
		m.aggregator.aggregate(tree)
	} else {
		m.sender.Offer(tree)
	}
}

func (m *Manager) isHitSample() bool {
	var sampleRatio float64

	if sampleRatio >= 1.0 {
		return true
	} else if sampleRatio < 1e-9 {
		return false
	} else {
		count := atomic.AddUint64(&m.sampleCount, 1)
		return count%(uint64(1/sampleRatio)) == 0
	}
}
