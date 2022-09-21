package analyzer

import "github.com/Orlion/cat-agent/cat/message"

type TransactionData struct {
	mtype, name string
	count, fail int
	sum         int64
	durations   map[int]int
}

type EventAggregator struct {
}

func (ea *EventAggregator) logEvent(event *message.Event) {

}
