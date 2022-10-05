package cat

import (
	"context"
	"fmt"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
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
	datas map[string]*eventData
	ch    chan *message.Event
}

func newEventAggregator() *EventAggregator {
	return &EventAggregator{
		datas: make(map[string]*eventData),
		ch:    make(chan *message.Event, config.EventAggregatorChannelSize),
	}
}

func (ea *EventAggregator) run(ctx context.Context) {
	log.Info("event aggregator running...")
	ticker := time.NewTicker(config.EventAggregatorTickerDuration)

Loop:
	for {
		select {
		case event := <-ea.ch:
			ea.getOrDefault(event).add(event)
		case <-ticker.C:
			ea.flush()
		case <-ctx.Done():
			break Loop
		}
	}

	ticker.Stop()

	close(ea.ch)

	for event := range ea.ch {
		ea.getOrDefault(event).add(event)
	}

	ea.flush()

	log.Info("event aggregator exit")
}

func (ea *EventAggregator) logEvent(event *message.Event) {
	select {
	case ea.ch <- event:
	default:
		log.Warnf("event aggregatro's ch is full, event: %s,%s  has been discarded", event.GetType(), event.GetName())
	}
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

	trans := message.NewTransaction(config.TypeSystem, config.NameEventAggregator, message.SUCCESS, "", 0, nil, 0)
	for _, data := range ea.datas {
		event := message.NewEvent(data.t, data.name, message.SUCCESS, fmt.Sprintf("%c%d%c%d", config.BatchFlag, data.count, config.BatchSplit, data.fail), 0)
		trans.AddChild(event)
	}

	tree := message.NewMessageTree()
	tree.SetMessage(trans)
	tree.SetMessageId(GetNextId())
	tree.SetParentMessageId("")
	tree.SetRootMessageId("")
	tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
	tree.SetThreadId(config.ThreadIdCatAgent)
	tree.SetThreadName(config.ThreadNameCatAgent)
	tree.SetDiscard(false)

	Send(tree)
}
