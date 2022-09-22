package manager

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Orlion/cat-agent/cat/message"
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

	millis := transaction.GetDurationInMillis()
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

	buf.WriteRune(batchFlag)
	buf.WriteString(strconv.Itoa(td.count))
	buf.WriteRune(batchSplit)
	buf.WriteString(strconv.Itoa(td.fail))
	buf.WriteRune(batchSplit)
	buf.WriteString(strconv.FormatUint(uint64(td.sum), 10))
	buf.WriteRune(batchSplit)

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

	buf.WriteRune(batchSplit)
	return buf.String()
}

type TransactionAggregator struct {
	manager *Manager
	datas   map[string]*transactionData
}

func newTransactionAggregator(manager *Manager) *TransactionAggregator {
	return &TransactionAggregator{
		manager: manager,
	}
}

func (ta *TransactionAggregator) logTransaction(transaction *message.Transaction) {
	ta.getOrDefault(transaction).add(transaction)
}

func (ta *TransactionAggregator) getOrDefault(transaction *message.Transaction) *transactionData {
	key := fmt.Sprintf("%s,%s", transaction.GetType(), transaction.GetName())

	data, exists := ta.datas[key]
	if !exists {
		data := &transactionData{
			t:         transaction.GetType(),
			name:      transaction.GetName(),
			count:     0,
			fail:      0,
			sum:       0,
			durations: make(map[int]int),
		}

		ta.datas[key] = data
	}

	return data
}

func (ta *TransactionAggregator) flush() {
	if len(ta.datas) == 0 {
		return
	}

	trans := message.NewTransaction(typeSystem, nameTransactionAggregator, message.SUCCESS, "", 0)

	for _, data := range ta.datas {
		trans := message.NewTransaction(data.t, data.name, message.SUCCESS, data.encode(), 0)
		trans.AddChild(trans)
	}

	ta.manager.Flush(trans, "todo", "", "", ThreadGroupNameCatAgent, ThreadIdCatAgent, ThreadNameCatAgent, false)
}
