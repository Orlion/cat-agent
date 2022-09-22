package server

import (
	"bufio"
	"net"
)

type conn struct {
	server     *Server
	rwc        net.Conn
	remoteAddr string
	bufr       *bufio.Reader
}

func (c *conn) serve() {
	c.remoteAddr = c.rwc.RemoteAddr().String()

	for {
		if c.server.shuttingDown() {
			break
		}

		req, err := c.readRequest()
		if err != nil {
			return
		}

		if handler, exists := c.server.handlers[req.Cmd]; exists {
			c.sendResponse(handler(req))
		}
	}

	c.close()
}

func (c *conn) close() {
	c.rwc.Close()
	c.bufr = nil
}
