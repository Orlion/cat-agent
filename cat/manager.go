package cat

import (
	"sync/atomic"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/cat/sender"
	"github.com/Orlion/cat-agent/log"
)

type Manager struct {
	aggregator  *LocalAggregator
	sender      sender.Sender
	sampleCount uint64
}

func newManager() *Manager {
	manager := &Manager{
		sender:     sender.NewTcpSender(),
		aggregator: newLocalAggregator(),
	}

	return manager
}

func (m *Manager) run() {
	log.Info("manager running...")
	m.sender.Run()
	m.aggregator.run()
}

func (m *Manager) shutdown() {
	log.Info("manager shutdown...")
	m.aggregator.shutdown()
	m.sender.Shutdown()
}

func (m *Manager) send(tree *message.MessageTree) {
	if tree.CanDiscard() && !m.hitSample() {
		m.aggregator.aggregate(tree)
	} else {
		m.sender.Offer(tree)
	}
}

func (m *Manager) hitSample() bool {
	sampleRatio := config.GetInstance().GetSample()

	if sampleRatio >= 1.0 {
		return true
	} else if sampleRatio < 1e-9 {
		return false
	} else {
		count := atomic.AddUint64(&m.sampleCount, 1)
		return count%(uint64(1/sampleRatio)) == 0
	}
}
