package cat

import (
	"sync/atomic"

	"github.com/Orlion/cat-agent/cat/analyzer"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/cat/sender"
)

type Manager struct {
	aggregator  *analyzer.LocalAggregator
	sender      *sender.TcpSender
	sampleCount uint64
}

func newManager() *Manager {
	return &Manager{
		aggregator: analyzer.NewLocalAggregator(),
		sender:     sender.NewTcpSender(),
	}
}

func (m *Manager) flush(domain string, hostname string, ipAddress string, msg message.Message, messageId string, parentMessageId string, rootMessageId string, threadGroupName string, threadId string, threadName string) {
	tree := message.NewMessageTree(domain, hostname, ipAddress, msg, messageId, parentMessageId, rootMessageId, threadGroupName, threadId, threadName)

	if tree.GetMessage().GetStatus() == message.SUCCESS && m.isHitSample() {
		m.aggregator.Aggregate(tree)
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
