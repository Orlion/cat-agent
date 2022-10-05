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

	TcpSenderHighQueueSize   = 50000
	TcpSenderNormalQueueSize = 50000

	TcpSenderNormalQueueConsumerNum      = 10
	TcpSenderHighQueueConsumerNum        = 10
	TcpSenderQueueConsumerTickerDuration = 1000 * time.Millisecond
	TcpSenderQueueConsumerBufSize        = 150

	EventAggregatorTickerDuration       = 3 * time.Second
	TransactionAggregatorTickerDuration = 3 * time.Second
	EventAggregatorChannelSize          = 1000
	TransactionAggregatorChannelSize    = 1000

	RouterUpdateDuration = 60 * time.Second

	BinaryProtocol = "NT1"
)
