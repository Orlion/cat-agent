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
	defer func() {
		c.close()
		c.server.decrConnNum()
		if err := recover(); err != nil {
			log.Errorf("conn serve panic, err: %v", err)
		}
	}()

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
			break
		}

		if handler, exists := c.server.handlers[req.Cmd]; exists {
			status, payload := handler(req)
			if req.Cmd != CmdSendMessage {
				err = c.sendResponse(status, payload)
				if err != nil {
					log.Errorf("conn send response error: %s", err)
					break
				}
			}

		} else {
			err = c.sendResponse(StatusNotFoundCmd, nil)
			if err != nil {
				log.Errorf("conn send response error: %s", err)
				break
			}
		}
	}
}

func (c *conn) close() {
	c.rwc.Close()
	c.bufr = nil
}
