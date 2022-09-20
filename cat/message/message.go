package message

type Message interface {
	AddData(k string, v string)
	SetData(data string)
	Complete()
	GetData() map[string]string
	GetName() string
	GetStatus() string
	GetTimestamp() int
	GetType() string
	IsCompleted() bool
	IsSuccess() bool
	SetStatus(status string)
	SetSuccessStatus()
	SetTimestamp(timestamp int)
}
