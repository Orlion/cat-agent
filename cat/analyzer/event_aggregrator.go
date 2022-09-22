package analyzer

import (
	"fmt"

	"github.com/Orlion/cat-agent/cat"
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
	datas map[string]map[string]*eventData
}

func (ea *EventAggregator) logEvent(domain string, event *message.Event) {

}

func (ea *EventAggregator) getOrDefault(domain string, event *message.Event) *eventData {
	var data *eventData

	key := fmt.Sprintf("%s,%s", event.GetType(), event.GetName())

	domainDatas, exists := ea.datas[domain]
	if !exists {
		data = &eventData{
			t:     event.GetType(),
			name:  event.GetName(),
			count: 0,
			fail:  0,
		}

		ea.datas[domain] = map[string]*eventData{
			key: data,
		}
	} else {
		data, exists = domainDatas[key]
		if !exists {
			data = &eventData{
				t:     event.GetType(),
				name:  event.GetName(),
				count: 0,
				fail:  0,
			}
		}
	}

	return data
}

func (ea *EventAggregator) flush() {
	if len(ea.datas) == 0 {
		return
	}

	for domain, domainDatas := range ea.datas {
		trans := message.NewTransaction(typeSystem, nameEventAggregator, message.SUCCESS, "", 0)

		for _, data := range domainDatas {
			event := message.NewEvent(data.t, data.name, message.SUCCESS, fmt.Sprintf("%c%d%c%d", batchFlag, data.count, batchSplit, data.fail), 0)
			trans.AddChild(event)
		}

		cat.Flush(domain, "todo", "todo", trans, "todo", "", "", "todo", "todo", "todo", false)
	}
}
