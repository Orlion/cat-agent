package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/dsx"
	"github.com/Orlion/cat-agent/server"
)

var (
	errBodyEof = errors.New("body eof")
	errBodyEnd = errors.New("body end")
)

const (
	Tab   byte = '\t'
	Lf    byte = '\n'
	TypeA byte = 'A'
	Typet byte = 't'
	TypeT byte = 'T'
	TypeE byte = 'E'
)

func SendMessage(req *server.Request) (status server.Status, payload []byte) {
	// read header
	r := &messageTreeReader{
		len:  len(req.Body),
		body: req.Body,
		tree: message.NewMessageTree(),
	}

	err := r.readHeader()
	if err != nil {
		log.Errorf("send message handler read header error: %s", err.Error())
		status = server.StatusMsgReadHeaderErr
		return
	}

	log.Debugf("read header, domain: %s, threadGroupName: %s, threadId: %s, threadName: %s, messageId: %s, parentMessageId: %s, rootMessageId: %s", r.tree.GetDomain(), r.tree.GetThreadGroupName(), r.tree.GetThreadId(), r.tree.GetThreadName(), r.tree.GetMessageId(), r.tree.GetParentMessageId(), r.tree.GetRootMessageId())

	err = r.readMessage()
	if err != nil {
		log.Errorf("send message handler read message error: %s", err.Error())
		status = server.StatusMsgReadMessageErr
		return
	}

	cat.Send(r.tree)

	return
}

type messageTreeReader struct {
	i    int
	len  int
	body []byte
	tree *message.MessageTree
}

func (r *messageTreeReader) readHeader() error {
	domain, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetDomain(domain)

	threadGroupName, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadGroupName(threadGroupName)
	threadId, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadId(threadId)

	threadName, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadName(threadName)

	messageId, err := r.readElement()
	if err != nil {
		return err
	}
	if len(messageId) > 0 {
		r.tree.SetMessageId(messageId)
	} else {
		r.tree.SetMessageId(cat.CreateMessageId(string(domain)))
	}

	parentMessageId, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetParentMessageId(parentMessageId)

	rootMessageId, err := r.readElement()
	if err == errBodyEnd {
		err = nil
	}
	r.tree.SetRootMessageId(rootMessageId)

	return nil
}

func (r *messageTreeReader) readMessage() (err error) {
	var (
		root  message.Message
		msg   message.Message
		stack = dsx.NewTransactionStack()
		t     byte
	)

Loop:
	for {
		t, msg, err = r.readMessageLine()
		if err != nil {
			break Loop
		}
		switch t {
		case Typet:
			if trans := stack.Peek(); trans != nil {
				stack.Peek().AddChild(msg)
			}
			stack.Push(msg.(*message.Transaction))
		case TypeT:
			if trans := stack.Pop(); trans == nil {
				err = errors.New("transaction are not a pair")
				break Loop
			} else if trans.GetType() != msg.GetType() || trans.GetName() != msg.GetName() {
				err = errors.New("transaction are not a pair")
				break Loop
			} else {
				root = trans
			}
		case TypeA:
			fallthrough
		case TypeE:
			if trans := stack.Peek(); trans != nil {
				stack.Peek().AddChild(msg)
			}
			root = msg
		}
		if msg.GetStatus() != message.SUCCESS {
			r.tree.SetDiscard(false)
		}
	}

	if err == nil {
		err = errors.New("body not eof")
	} else if err == errBodyEof && stack.IsEmpty() {
		err = nil
		r.tree.SetMessage(root)
	}

	return
}

func (r *messageTreeReader) readMessageLine() (t byte, msg message.Message, err error) {
	tBytes, err := r.readElement()
	if err != nil {
		return
	}

	if len(tBytes) != 1 {
		err = fmt.Errorf("unknown type: %s", string(t))
		return
	}

	t = tBytes[0]

	mtype, err := r.readElement()
	if err != nil {
		return
	}

	name, err := r.readElement()
	if err != nil {
		return
	}

	status, err := r.readElement()
	if err != nil {
		return
	}

	timestampInMillis, err := r.readElement()
	if err != nil {
		return
	}

	timestampInMillisInt64, _ := strconv.ParseInt(string(timestampInMillis), 10, 64)

	durationInMicros, err := r.readElement()
	if err != nil {
		return
	}

	data, err := r.readElement()
	if err == errBodyEnd {
		err = nil
	}

	switch t {
	case Typet:
		fallthrough
	case TypeT:
		fallthrough
	case TypeA:
		durationInMicrosInt64, _ := strconv.ParseInt(string(durationInMicros), 10, 64)
		msg = message.NewTransaction(string(mtype), string(name), string(status), string(data), timestampInMillisInt64, nil, durationInMicrosInt64)
	case TypeE:
		msg = message.NewEvent(string(mtype), string(name), string(status), string(data), timestampInMillisInt64)
	default:
		err = fmt.Errorf("unknown type: %s", string(t))
	}

	return
}

func (r *messageTreeReader) readElement() (b []byte, err error) {
	b = make([]byte, 0)

	for {
		if r.i >= r.len {
			err = errBodyEof
			return
		}

		if r.body[r.i] == Lf {
			err = errBodyEnd
			break
		}
		if r.body[r.i] == Tab {
			break
		}
		b = append(b, r.body[r.i])
		r.i++
	}

	r.i++

	return
}
