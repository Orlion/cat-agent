package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/dsx"
	"github.com/Orlion/cat-agent/server"
	"gitlab-team.smzdm.com/smzdm/zdm-go-cat/message"
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
	r := &messageReader{
		len:  len(req.Body),
		body: req.Body,
	}

	header, err := r.readHeader()
	if err != nil {
		log.Errorf("send message handler read header error: %s", err.Error())
		status = server.StatusMsgReadHeaderErr
		return
	}

	return
}

type messageHeader struct {
	domain          []byte
	threadGroupName []byte
	threadId        []byte
	threadName      []byte
	messageId       []byte
	parentMessageId []byte
	rootMessageId   []byte
}

type messageReader struct {
	i    int
	len  int
	body []byte
}

func (r *messageReader) readHeader() (header *messageHeader, err error) {
	header = new(messageHeader)

	header.domain, err = r.readElement()
	if err != nil {
		return
	}

	header.threadGroupName, err = r.readElement()
	if err != nil {
		return
	}

	header.threadId, err = r.readElement()
	if err != nil {
		return
	}

	header.threadName, err = r.readElement()
	if err != nil {
		return
	}

	header.messageId, err = r.readElement()
	if err != nil {
		return
	}

	header.parentMessageId, err = r.readElement()
	if err != nil {
		return
	}

	header.rootMessageId, err = r.readElement()
	if err == errBodyEnd {
		err = nil
	}

	return
}

func (r *messageReader) readMessage() (msg message.Message, err error) {
	var (
		t     byte
		stack = dsx.NewTransactionStack()
	)
Loop:
	for {
		t, msg, err = r.readMessageLine()
		if err == errBodyEnd {
			err = nil
		}
		if err != nil {
			return
		}

		switch t {
		case Typet:
			stack.Push(msg)
		case TypeT:

		case TypeA:
		case TypeE:
			break Loop
		}
	}

	return
}

func (r *messageReader) readMessageLine() (t byte, msg message.Message, err error) {
	tBytes, err := r.readElement()
	if err != nil {
		return
	}

	if len(tBytes) != 1 {
		err = fmt.Errorf("unknown type: %s", string(t))
		return
	}

	t = tBytes[0]

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
	case TypeT:
		rawDurationInMicrosInt64, _ := strconv.ParseInt(string(rawDurationInMicros), 10, 64)
		msg = message.NewTransaction(string(mtype), string(name), string(status), string(data), timestampInMillisInt64, nil, rawDurationInMicrosInt64)
	case TypeE:
		msg = message.NewEvent(string(mtype), string(name), string(status), string(data), timestampInMillisInt64)
	default:
		err = fmt.Errorf("unknown type: %s", string(t))
	}

	return
}

func (r *messageReader) readElement() (b []byte, err error) {
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
