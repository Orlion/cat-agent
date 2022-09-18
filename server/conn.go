package server

import (
	"net"
)

type conn struct {
	server     *Server
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) serve() {
	c.remoteAddr = c.rwc.RemoteAddr().String()

	for {
		req, err := c.readRequest()
		if err != nil {
			return
		}

		if handler, exists := c.server.handlers[req.Action]; exists {
			c.sendResponse(handler(req))
		}
	}
}
