package server

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var ErrServerClosed = errors.New("server: Server closed")

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	mu           sync.Mutex
	doneChan     chan struct{}
}

func (srv *Server) ListenAndServe() error {
	network := "tcp"

	if strings.HasPrefix(srv.Addr, "unix:") {
		os.Remove(srv.Addr)
		network = "unix"
	}

	ln, err := net.Listen(network, srv.Addr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

func (srv *Server) Serve(ln net.Listener) error {
	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		rw, err := ln.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		c := srv.newConn(rw)
		go c.serve()
	}
}

func (src *Server) Shutdown(ctx context.Context) {

}

func (s *Server) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getDoneChanLocked()
}

func (s *Server) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}

	return c
}
