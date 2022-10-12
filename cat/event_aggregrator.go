package cat

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/timex"
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

type eventWithDomain struct {
	domain string
	event  *message.Event
}

type EventAggregator struct {
	datas map[string]map[string]*eventData
	ch    chan *eventWithDomain
}

func newEventAggregator() *EventAggregator {
	return &EventAggregator{
		datas: make(map[string]map[string]*eventData),
		ch:    make(chan *eventWithDomain, config.EventAggregatorChannelSize),
	}
}

func (ea *EventAggregator) run(ctx context.Context) {
	log.Info("event aggregator running...")
	ticker := time.NewTicker(config.EventAggregatorTickerDuration)

Loop:
	for {
		select {
		case eventWithDomain := <-ea.ch:
			ea.getOrDefault(eventWithDomain).add(eventWithDomain.event)
		case <-ticker.C:
			ea.flush()
		case <-ctx.Done():
			break Loop
		}
	}

	ticker.Stop()

	close(ea.ch)

	for eventWithDomain := range ea.ch {
		ea.getOrDefault(eventWithDomain).add(eventWithDomain.event)
	}

	ea.flush()

	log.Info("event aggregator exit")
}

func (ea *EventAggregator) logEvent(domain string, event *message.Event) {
	select {
	case ea.ch <- &eventWithDomain{domain, event}:
	default:
		log.Warnf("event aggregatro's ch is full, event: %s,%s  has been discarded", event.GetType(), event.GetName())
	}
}

func (ea *EventAggregator) getOrDefault(eventWithDomain *eventWithDomain) (data *eventData) {
	key := fmt.Sprintf("%s,%s", eventWithDomain.event.GetType(), eventWithDomain.event.GetName())

	if domainDatas, exists := ea.datas[eventWithDomain.domain]; exists {
		if data, exists = domainDatas[key]; !exists {
			data = &eventData{
				t:     eventWithDomain.event.GetType(),
				name:  eventWithDomain.event.GetName(),
				count: 0,
				fail:  0,
			}

			domainDatas[eventWithDomain.domain] = data
		}
	} else {
		data = &eventData{
			t:     eventWithDomain.event.GetType(),
			name:  eventWithDomain.event.GetName(),
			count: 0,
			fail:  0,
		}

		ea.datas[eventWithDomain.domain] = map[string]*eventData{key: data}
	}

	return data
}

func (ea *EventAggregator) flush() {
	if len(ea.datas) == 0 {
		return
	}

	for domain, domainDatas := range ea.datas {
		trans := message.NewTransaction(config.TypeSystem, config.NameEventAggregator, message.SUCCESS, "", timex.NowUnixMillis(), nil, 0)

		for _, data := range domainDatas {
			child := message.NewEvent(data.t, data.name, message.SUCCESS, fmt.Sprintf("%c%d%c%d", config.BatchFlag, data.count, config.BatchSplit, data.fail), timex.NowUnixMillis())
			trans.AddChild(child)
		}

		tree := message.NewMessageTree()
		tree.SetMessage(trans)
		tree.SetDomain([]byte(domain))
		messageId := CreateMessageId(domain)
		tree.SetMessageId(messageId)
		tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
		tree.SetThreadId([]byte(strconv.Itoa(os.Getpid())))
		tree.SetThreadName(config.ThreadNameCatAgent)
		tree.SetDiscard(false)

		Send(tree)

		log.Debugf("event aggregator flush, messageId: %s, ", messageId)
	}

	ea.datas = make(map[string]map[string]*eventData)
}
