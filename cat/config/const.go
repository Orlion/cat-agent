package config

import "time"

const (
	DefaultHostname = "UnknownHost"
	DefaultEnv      = "dev"
	DefaultIp       = "127.0.0.1"
	DefaultIpHex    = "7f000001"

	TypeSystem    = "System"
	TypeHeartbeat = "Heartbeat"

	NameReboot                = "Reboot"
	NameTransactionAggregator = "TransactionAggregator"
	NameEventAggregator       = "EventAggregator"
	NameStatus                = "Status"
	NameStatusExtensionPrefix = "StatusExtension-"

	BatchFlag  = '@'
	BatchSplit = ';'

	TcpSenderHighQueueSize   = 50000
	TcpSenderNormalQueueSize = 50000

	DefaultTcpSenderNormalQueueConsumerNum = 10
	DefaultTcpSenderHighQueueConsumerNum   = 10
	TcpSenderQueueConsumerTickerDuration   = 1000 * time.Millisecond
	TcpSenderQueueConsumerBufSize          = 150

	EventAggregatorTickerDuration       = 3 * time.Second
	TransactionAggregatorTickerDuration = 3 * time.Second
	EventAggregatorChannelSize          = 1000
	TransactionAggregatorChannelSize    = 1000

	RouterUpdateDuration = 60 * time.Second
)

var (
	BinaryProtocol          = []byte("NT1")
	ThreadNameCatAgent      = []byte("cat-agent")
	ThreadGroupNameCatAgent = []byte("cat-agent-group")
)
