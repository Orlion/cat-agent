package server

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/Orlion/cat-agent/log"
)

type Cmd uint32

const (
	CmdCreateMessageId Cmd = iota + 1
	CmdSendMessage
)

const ReqHeaderLen = 8

type Request struct {
	Cmd    Cmd
	Length uint32
	Body   []byte
}

func (c *conn) readRequest() (req *Request, err error) {
	if c.server.ReadTimeout != 0 {
		err = c.rwc.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
		if err != nil {
			return nil, err
		}
	}

	// read header
	buf := make([]byte, 8)
	_, err = io.ReadFull(c.bufr, buf)
	if err != nil {
		return
	}

	req = new(Request)
	req.Cmd = Cmd(binary.BigEndian.Uint32(buf[0:4]))
	req.Length = binary.BigEndian.Uint32(buf[4:8])

	log.Debugf("recv request from %s, cmd: %d, length: %d", c.rwc.RemoteAddr().String(), req.Cmd, req.Length)

	if req.Length > ReqHeaderLen {
		// read body
		req.Body = make([]byte, req.Length-ReqHeaderLen)
		_, err = io.ReadFull(c.bufr, req.Body)
		if err != nil {
			return
		}
	}

	return
}
