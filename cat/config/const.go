package config

import "time"

const (
	DefaultHostname = "UnknownHost"
	DefaultEnv      = "dev"
	DefaultIp       = "127.0.0.1"
	DefaultIpHex    = "7f000001"

	TypeSystem = "System"

	NameReboot = "Reboot"

	NameTransactionAggregator = "TransactionAggregator"
	NameEventAggregator       = "EventAggregator"

	BatchFlag  = '@'
	BatchSplit = ';'

	ThreadIdCatAgent        = "0"
	ThreadNameCatAgent      = "cat-agent"
	ThreadGroupNameCatAgent = "cat-agent-group"

	HighPriorityQueueSize   = 50000
	NormalPriorityQueueSize = 50000

	NormalQueueConsumerNum      = 10
	HighQueueConsumerNum        = 10
	QueueConsumerTickerDuration = 100 * time.Millisecond
	QueueConsumerBufSize        = 150

	EventAggregatorTickerDuration       = time.Second * 3
	TransactionAggregatorTickerDuration = time.Second * 3
)
