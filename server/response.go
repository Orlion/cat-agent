package server

import (
	"encoding/binary"
	"time"

	"github.com/Orlion/cat-agent/log"
)

type Status uint32

const (
	StatusOk Status = iota
	StatusMsgReadHeaderErr
	StatusMsgReadMessageErr
	StatusNotFoundCmd
	StatusBadDomain
)

const RespHeaderLen = 8

type response struct {
	status  Status
	length  uint32
	payload []byte
}

func (c *conn) sendResponse(status Status, payload []byte) (err error) {
	if c.server.WriteTimeout != 0 {
		err = c.rwc.SetWriteDeadline(time.Now().Add(c.server.ReadTimeout))
		if err != nil {
			return
		}
	}

	length := RespHeaderLen + uint32(len(payload))
	b := make([]byte, length)
	binary.BigEndian.PutUint32(b, uint32(status))
	binary.BigEndian.PutUint32(b[4:8], length)
	copy(b[RespHeaderLen:], payload)

	log.Debugf("send response to %s, status: %d, length: %d", c.rwc.RemoteAddr().String(), status, length)

	var n int

	for {
		n, err = c.rwc.Write(b)
		if err != nil {
			return
		}

		if len(b) <= n {
			break
		}

		b = b[n:]
	}

	return
}
