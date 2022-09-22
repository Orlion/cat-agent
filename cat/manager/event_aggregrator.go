package manager

import (
	"fmt"

	"github.com/Orlion/cat-agent/cat/message"
)

type eventData struct {
	t, name     string
	count, fail int
}

func (ed *eventData) add(event *message.Event) {
	ed.count++
	if event.GetStatus() != message.SUCCESS {
		ed.fail++
	}
}

type EventAggregator struct {
	manager *Manager
	datas   map[string]*eventData
}

func newEventAggregator(manager *Manager) *EventAggregator {
	return &EventAggregator{
		manager: manager,
	}
}

func (ea *EventAggregator) logEvent(event *message.Event) {

}

func (ea *EventAggregator) getOrDefault(event *message.Event) *eventData {
	key := fmt.Sprintf("%s,%s", event.GetType(), event.GetName())

	data, exists := ea.datas[key]
	if !exists {
		data = &eventData{
			t:     event.GetType(),
			name:  event.GetName(),
			count: 0,
			fail:  0,
		}
	}

	return data
}

func (ea *EventAggregator) flush() {
	if len(ea.datas) == 0 {
		return
	}

	trans := message.NewTransaction(typeSystem, nameEventAggregator, message.SUCCESS, "", 0, nil, 0)
	for _, data := range ea.datas {
		event := message.NewEvent(data.t, data.name, message.SUCCESS, fmt.Sprintf("%c%d%c%d", batchFlag, data.count, batchSplit, data.fail), 0)
		trans.AddChild(event)
	}

	ea.manager.Flush(trans, "todo", "", "", ThreadGroupNameCatAgent, ThreadIdCatAgent, ThreadNameCatAgent, false)
}
