package server

import (
	"net"
	"time"
)

type conn struct {
	server     *Server
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) serve() {
	c.remoteAddr = c.rwc.RemoteAddr().String()

	for {
		w, err := c.readRequest()
	}
}

func (c *conn) readRequest() (w *response, err error) {
	var (
		reqDeadline time.Time
	)

	t0 := time.Now()
	if d := c.server.ReadTimeout; d != 0 {
		reqDeadline = t0.Add(d)
	}
	c.rwc.SetReadDeadline(reqDeadline)
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	req, err := readRequest()
	if err != nil {
		return
	}

	w := &response{}
}
