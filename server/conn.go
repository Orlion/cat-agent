package server

import (
	"bufio"
	"errors"
	"io"
	"net"

	"github.com/Orlion/cat-agent/log"
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
			if errors.Is(err, io.EOF) {
				log.Infof("conn from %s closed", c.remoteAddr)
			} else {
				log.Errorf("conn read request from %s error: %s", c.remoteAddr, err.Error())
			}
			return
		}

		if handler, exists := c.server.handlers[req.Cmd]; exists {
			err = c.sendResponse(handler(req))
			if err != nil {
				log.Errorf("conn send response error: %s", err)
				return
			}
		} else {
			err = c.sendResponse(StatusNotFoundCmd, nil)
			if err != nil {
				log.Errorf("conn send response error: %s", err)
				return
			}
		}
	}

	c.close()
}

func (c *conn) close() {
	c.rwc.Close()
	c.bufr = nil
}
