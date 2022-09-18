package server

import (
	"encoding/binary"
	"time"
)

type Status uint32

const (
	StatusOk Status = iota
)

const RespHeaderLen = 8

type response struct {
	status Status
	length uint32
	body   []byte
}

func (c *conn) sendResponse(status Status, body []byte) error {
	if c.server.WriteTimeout != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(c.server.ReadTimeout))
	}

	resp := &response{
		status: status,
		length: RespHeaderLen + uint32(len(body)),
		body:   body,
	}

	b := make([]byte, 0, resp.length)
	binary.BigEndian.PutUint32(b, uint32(resp.status))
	binary.BigEndian.PutUint32(b, uint32(resp.length))
	b = append(b, resp.body...)

	c.rwc.Write(b)

	return nil
}
