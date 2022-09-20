package message

import (
	"github.com/Orlion/cat-agent/cat/analyzer"
	"github.com/Orlion/cat-agent/cat/sender"
)

type Manager struct {
	tree       *MessageTree
	aggregator *analyzer.LocalAggregator
	sender     *sender.Sender
}

func (m *Manager) flush() {
	if m.tree.canDiscard() && m.isHitSample() {
		m.aggregator.Aggregate(m.tree)
	} else {
		m.sender.Offer(m.tree)
	}
}

func (m *Manager) isHitSample() bool {
	return false
}
