package cat

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
)

type transactionData struct {
	t, name     string
	count, fail int
	sum         int64
	durations   map[int]int
}

func (td *transactionData) add(transaction *message.Transaction) {
	td.count++

	if transaction.GetStatus() != message.SUCCESS {
		td.fail++
	}

	millis := transaction.GetDurationInMicros() / 1000
	td.sum += millis

	duration := computeDuration(int(millis))
	if _, ok := td.durations[duration]; ok {
		td.durations[duration]++
	} else {
		td.durations[duration] = 1
	}
}

func (td *transactionData) encode() string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteRune(config.BatchFlag)
	buf.WriteString(strconv.Itoa(td.count))
	buf.WriteRune(config.BatchSplit)
	buf.WriteString(strconv.Itoa(td.fail))
	buf.WriteRune(config.BatchSplit)
	buf.WriteString(strconv.FormatUint(uint64(td.sum), 10))
	buf.WriteRune(config.BatchSplit)

	i := 0
	for k, v := range td.durations {
		if i > 0 {
			buf.WriteRune('|')
		}
		buf.WriteString(strconv.Itoa(k))
		buf.WriteRune(',')
		buf.WriteString(strconv.Itoa(v))
		i++
	}

	buf.WriteRune(config.BatchSplit)
	return buf.String()
}

type transactionWithDomain struct {
	domain      string
	transaction *message.Transaction
}

type TransactionAggregator struct {
	datas map[string]map[string]*transactionData
	ch    chan *transactionWithDomain
}

func newTransactionAggregator() *TransactionAggregator {
	return &TransactionAggregator{
		datas: make(map[string]map[string]*transactionData),
		ch:    make(chan *transactionWithDomain, config.TransactionAggregatorChannelSize),
	}
}

func (ta *TransactionAggregator) run(ctx context.Context) {
	log.Info("transaction aggregator running...")
	ticker := time.NewTicker(config.TransactionAggregatorTickerDuration)

Loop:
	for {
		select {
		case transWithDomain := <-ta.ch:
			ta.getOrDefault(transWithDomain.domain, transWithDomain.transaction).add(transWithDomain.transaction)
		case <-ticker.C:
			ta.flush()
		case <-ctx.Done():
			break Loop
		}
	}

	ticker.Stop()

	close(ta.ch)

	for transWithDomain := range ta.ch {
		ta.getOrDefault(transWithDomain.domain, transWithDomain.transaction).add(transWithDomain.transaction)
	}

	ta.flush()
	log.Info("transaction aggregator exit")
}

func (ta *TransactionAggregator) logTransaction(domain string, transaction *message.Transaction) {
	select {
	case ta.ch <- &transactionWithDomain{domain, transaction}:
	default:
		log.Warnf("transaction aggregatro's ch is full, transaction: %s,%s has been discarded", transaction.GetType(), transaction.GetName())
	}
}

func (ta *TransactionAggregator) getOrDefault(domain string, transaction *message.Transaction) (data *transactionData) {
	key := fmt.Sprintf("%s,%s", transaction.GetType(), transaction.GetName())

	domainData, exists := ta.datas[domain]
	if exists {
		data, exists = domainData[key]
		if !exists {
			data = &transactionData{
				t:         transaction.GetType(),
				name:      transaction.GetName(),
				count:     0,
				fail:      0,
				sum:       0,
				durations: make(map[int]int),
			}

			domainData[key] = data
		}
	} else {
		data = &transactionData{
			t:         transaction.GetType(),
			name:      transaction.GetName(),
			count:     0,
			fail:      0,
			sum:       0,
			durations: make(map[int]int),
		}

		ta.datas[domain] = map[string]*transactionData{key: data}
	}

	return
}

func (ta *TransactionAggregator) flush() {
	if len(ta.datas) == 0 {
		return
	}

	for domain, datas := range ta.datas {
		trans := message.NewTransaction(config.TypeSystem, config.NameTransactionAggregator, message.SUCCESS, "", 0, nil, 0)

		for _, data := range datas {
			trans := message.NewTransaction(data.t, data.name, message.SUCCESS, data.encode(), 0, nil, 0)
			trans.AddChild(trans)
		}

		tree := message.NewMessageTree()
		tree.SetMessage(trans)
		tree.SetMessageId(GetNextId(domain))
		tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
		tree.SetThreadId(config.ThreadIdCatAgent)
		tree.SetThreadName(config.ThreadNameCatAgent)
		tree.SetDiscard(false)
		Send(tree)
	}

	ta.datas = make(map[string]map[string]*transactionData)

}
