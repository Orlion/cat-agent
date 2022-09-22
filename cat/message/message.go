package message

const SUCCESS = "0"

type Message interface {
	GetType() string
	GetName() string
	GetStatus() string
	GetData() string
	GetTimestamp() int64
	IsSuccess() bool
}

type baseMessage struct {
	t                 string
	name              string
	status            string
	data              string
	timestampInMillis int64
}

func newBaseMessage(t, name, status, data string, timestampInMillis int64) baseMessage {
	return baseMessage{
		t:                 t,
		name:              name,
		status:            status,
		data:              data,
		timestampInMillis: timestampInMillis,
	}
}

func (m *baseMessage) GetType() string {
	return m.t
}

func (m *baseMessage) GetName() string {
	return m.name
}

func (m *baseMessage) GetStatus() string {
	return m.status
}

func (m *baseMessage) GetData() string {
	return m.data
}

func (m *baseMessage) GetTimestamp() int64 {
	return m.timestampInMillis
}

func (m *baseMessage) IsSuccess() bool {
	return m.status == SUCCESS
}
