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
	Tab   = '\t'
	Lf    = '\n'
	TypeA = 'A'
	Typet = 't'
	TypeT = 'T'
	TypeE = 'E'
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
	i      int
	len    int
	body   []byte
	domain string
	tree   *message.MessageTree
}

func (r *messageTreeReader) readHeader() error {
	domain, err := r.readElement()
	if err != nil {
		return err
	}

	r.domain = string(domain)

	threadGroupName, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadGroupName(string(threadGroupName))

	threadId, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadId(string(threadId))

	threadName, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetThreadName(string(threadName))

	messageId, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetMessageId(string(messageId))

	parentMessageId, err := r.readElement()
	if err != nil {
		return err
	}
	r.tree.SetParentMessageId(string(parentMessageId))

	rootMessageId, err := r.readElement()
	if err == errBodyEnd {
		err = nil
	}
	r.tree.SetRootMessageId(string(rootMessageId))

	return nil
}

func (r *messageTreeReader) readMessage() (err error) {
	var (
		msg   message.Message
		stack = dsx.NewTransactionStack()
	)

	t, root, err := r.readMessageLine()
	if err != nil && err != errBodyEnd {
		return
	}

	for {
		switch t {
		case Typet:
			if trans := stack.Peek(); trans != nil {
				stack.Peek().AddChild(msg)
			}
			stack.Push(msg.(*message.Transaction))
		case TypeT:
			if trans := stack.Pop(); trans == nil {
				err = errors.New("transaction are not a pair")
				break
			}
		case TypeA:
			if trans := stack.Peek(); trans != nil {
				stack.Peek().AddChild(msg)
			}
		case TypeE:
			if trans := stack.Peek(); trans != nil {
				stack.Peek().AddChild(msg)
			}
		}

		t, msg, err = r.readMessageLine()
		if err == errBodyEnd {
			continue
		}
		if err != nil {
			break
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

	if t == TypeT {
		return
	}

	timestampInMillis, err := r.readElement()
	if err != nil {
		return
	}

	timestampInMillisInt64, _ := strconv.ParseInt(string(timestampInMillis), 10, 64)

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

	rawDurationInMicros, err := r.readElement()
	if err != nil {
		return
	}

	data, err := r.readElement()
	if err == errBodyEnd {
		err = nil
	}

	switch t {
	case Typet:
	case TypeA:
		rawDurationInMicrosInt64, _ := strconv.ParseInt(string(rawDurationInMicros), 10, 64)
		msg = message.NewTransaction(string(mtype), string(name), string(status), string(data), timestampInMillisInt64, nil, rawDurationInMicrosInt64)
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
			break
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

	return
}
