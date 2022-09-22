package server

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
)

var (
	ErrServerClosed = errors.New("server: Server closed")
)

type Handler func(req *Request) (status Status, payload []byte)

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	handlers map[Cmd]Handler

	inShutdown atomicx.Bool

	mu       sync.Mutex
	listener net.Listener
	doneChan chan struct{}
}

func NewServer(config *Config) *Server {
	withDefaultConf(config)
	return &Server{
		Addr:         config.Addr,
		ReadTimeout:  time.Duration(config.ReadTimeoutMillis) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeoutMillis) * time.Millisecond,
		handlers:     make(map[Cmd]Handler),
	}
}

func (srv *Server) Handle(cmd Cmd, handler Handler) {
	srv.handlers[cmd] = handler
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

	return srv.serve(ln)
}

func (srv *Server) serve(ln net.Listener) error {
	log.Infof("server listen on %s...", srv.Addr)

	var tempDelay time.Duration // how long to sleep on accept failure

	srv.listener = ln

	for {
		rw, err := ln.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Warnf("server accept temporary error: %s, tempDelay: %d", err, tempDelay)
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

func (srv *Server) Shutdown() error {
	srv.inShutdown.SetTrue()

	srv.mu.Lock()
	defer srv.mu.Unlock()

	lnerr := srv.listener.Close()
	srv.closeDoneChanLocked()

	return lnerr
}

func (srv *Server) shuttingDown() bool {
	return srv.inShutdown.Get()
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

func (s *Server) closeDoneChanLocked() {
	ch := s.getDoneChanLocked()
	select {
	case <-ch:
	default:
		close(ch)
	}
}

func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
		bufr:   bufio.NewReader(rwc),
	}

	return c
}
