package server

import (
	"bufio"
	"encoding/binary"
	"io"
	"time"
)

type Action uint32

const (
	ActionCreateMessageId Action = iota
	ActionSendMessage
)

const ReqHeaderLen = 8

type Request struct {
	Action Action
	Length uint32
	Body   []byte
}

func (c *conn) readRequest() (req *Request, err error) {
	if c.server.ReadTimeout != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
	}

	bufr := bufio.NewReader(c.rwc)

	// read header
	buf := make([]byte, 8)
	_, err = io.ReadFull(bufr, buf)
	if err != nil {
		return
	}

	req = new(Request)
	req.Action = Action(binary.BigEndian.Uint32(buf[0:4]))
	req.Length = binary.BigEndian.Uint32(buf[4:4])

	if req.Length > ReqHeaderLen {
		// read body
		req.Body = make([]byte, req.Length-8)
		_, err = io.ReadFull(bufr, req.Body)
		if err != nil {
			return
		}
	}

	return
}
