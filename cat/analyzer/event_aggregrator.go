package analyzer

type TransactionData struct {
	mtype, name string
	count, fail int
	sum         int64
	durations   map[int]int
}

type EventAggregator struct {
}
