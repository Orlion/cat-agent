package encoder

import (
	"bytes"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
)

type BinaryEncoder struct {
	buf  *bytes.Buffer
	tree *message.MessageTree
}

func NewBinaryEncoder() *BinaryEncoder {
	return &BinaryEncoder{
		buf: bytes.NewBuffer([]byte{}),
	}
}

func (e *BinaryEncoder) BufLen() int {
	return e.buf.Len()
}

func (e *BinaryEncoder) Bytes() []byte {
	return e.buf.Bytes()
}

func (e *BinaryEncoder) EncodeMessageTree(tree *message.MessageTree) (err error) {
	e.buf.Reset()

	e.tree = tree
	if err = e.encodeHeader(); err != nil {
		return
	}

	if err = e.encodeBody(); err != nil {
		return
	}

	return
}

func (e *BinaryEncoder) encodeHeader() (err error) {
	if _, err = e.buf.Write(config.BinaryProtocol); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetDomain()); err != nil {
		return
	}
	if err = e.writeString(config.GetInstance().GetHostname()); err != nil {
		return
	}
	if err = e.writeString(config.GetInstance().GetIp()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetThreadGroupName()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetThreadId()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetThreadName()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetMessageId()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetParentMessageId()); err != nil {
		return
	}
	if err = e.writeBytes(e.tree.GetRootMessageId()); err != nil {
		return
	}

	// sessionToken.
	if err = e.writeString(""); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) encodeBody() (err error) {
	return e.encodeMessage(e.tree.GetMessage())
}

func (e *BinaryEncoder) encodeMessage(m message.Message) (err error) {
	switch m.(type) {
	case *message.Transaction:
		return e.encodeTransaction(m.(*message.Transaction))
	case *message.Event:
		return e.encodeEvent(m.(*message.Event))
	case *message.Heartbeat:
		return e.encodeHeartbeat(m.(*message.Heartbeat))
	default:
		return
	}
}

func (e *BinaryEncoder) encodeTransaction(trans *message.Transaction) (err error) {
	if _, err = e.buf.WriteRune('t'); err != nil {
		return
	}

	if err = e.encodeMessageStart(trans); err != nil {
		return
	}

	for _, child := range trans.GetChildren() {
		if err = e.encodeMessage(child); err != nil {
			return
		}
	}

	if _, err = e.buf.WriteRune('T'); err != nil {
		return
	}
	if err = e.encodeMessageEnd(trans); err != nil {
		return
	}

	if err = e.writeI64(trans.GetDurationInMicros()); err != nil {
		return
	}

	return
}

func (e *BinaryEncoder) encodeEvent(event *message.Event) (err error) {
	return e.encodeMessageWithLeader(event, 'E')
}

func (e *BinaryEncoder) encodeHeartbeat(heartbeat *message.Heartbeat) (err error) {
	return e.encodeMessageWithLeader(heartbeat, 'H')
}

func (e *BinaryEncoder) encodeMessageWithLeader(m message.Message, leader rune) (err error) {
	if _, err = e.buf.WriteRune(leader); err != nil {
		return
	}
	if err = e.encodeMessageStart(m); err != nil {
		return
	}
	if err = e.encodeMessageEnd(m); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) encodeMessageStart(m message.Message) (err error) {
	if err = e.writeI64(m.GetTimestamp()); err != nil {
		return
	}
	if err = e.writeString(m.GetType()); err != nil {
		return
	}
	if err = e.writeString(m.GetName()); err != nil {
		return
	}

	return
}

func (e *BinaryEncoder) encodeMessageEnd(m message.Message) (err error) {
	if err = e.writeString(m.GetStatus()); err != nil {
		return
	}

	if m.GetData() == "" {
		if err = e.writeI64(0); err != nil {
			return
		}
	} else {
		if err = e.writeString(m.GetData()); err != nil {
			return
		}
	}
	return
}

func (e *BinaryEncoder) writeString(s string) (err error) {
	if err = e.writeI64(int64(len(s))); err != nil {
		return
	}
	if _, err = e.buf.WriteString(s); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) writeBytes(b []byte) (err error) {
	if err = e.writeI64(int64(len(b))); err != nil {
		return
	}
	if _, err = e.buf.Write(b); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) writeI64(i int64) (err error) {
	for {
		if i&^0x7F == 0 {
			err = e.buf.WriteByte(byte(i))
			return
		} else {
			if err = e.buf.WriteByte(byte(i&0x7F | 0x80)); err != nil {
				return
			}
			i >>= 7
		}
	}
}
