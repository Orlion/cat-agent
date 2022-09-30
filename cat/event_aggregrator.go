package cat

import (
	"fmt"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
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
	datas      map[string]*eventData
	ch         chan *message.Event
	inShutdown atomicx.Bool
	done       chan struct{}
}

func newEventAggregator() *EventAggregator {
	return &EventAggregator{
		datas: make(map[string]*eventData),
		ch:    make(chan *message.Event, config.EventAggregatorChannelSize),
		done:  make(chan struct{}),
	}
}

func (ea *EventAggregator) run() {
	log.Info("event aggregator running...")
	go func() {
		ticker := time.NewTicker(config.EventAggregatorTickerDuration)

		for !ea.inShutdown.Get() {
			select {
			case event := <-ea.ch:
				ea.getOrDefault(event).add(event)
			case <-ticker.C:
				ea.flush()
			}
		}

		ticker.Stop()

		close(ea.ch)

		for event := range ea.ch {
			ea.getOrDefault(event).add(event)
		}

		ea.flush()

		close(ea.done)
	}()
}

func (ea *EventAggregator) shutdown() {
	log.Info("event aggregator shutdown...")
	ea.inShutdown.SetTrue()
	<-ea.done
	log.Info("event aggregator exit")
}

func (ea *EventAggregator) logEvent(event *message.Event) {
	if ea.inShutdown.Get() {
		return
	}

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
	tree.SetMessageId(getNextId())
	tree.SetParentMessageId("")
	tree.SetRootMessageId("")
	tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
	tree.SetThreadId(config.ThreadIdCatAgent)
	tree.SetThreadName(config.ThreadNameCatAgent)
	tree.SetDiscard(false)

	catInstance.manager.flush(tree)
}
