package server

import (
	"bufio"
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
)

var (
	ErrServerClosed      = errors.New("server: Server closed")
	shutdownPollInterval = 500 * time.Millisecond
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
	connNum  int64
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

func (srv *Server) ListenAndServe() (err error) {
	network := "tcp"

	addr := srv.Addr
	if strings.HasPrefix(addr, "unix://") {
		addr = strings.TrimPrefix(addr, "unix://")
		os.Remove(addr)
		network = "unix"
	}

	srv.listener, err = net.Listen(network, addr)
	if err != nil {
		return err
	}

	go srv.serve()

	return nil
}

func (srv *Server) serve() error {
	log.Infof("server listen on %s...", srv.Addr)

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		rw, err := srv.listener.Accept()
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
		log.Debugf("server new conn from %s", rw.RemoteAddr().String())
		go func() {
			c.serve()
		}()
	}
}

func (srv *Server) Shutdown(ctx context.Context) error {
	log.Info("server shutdown...")

	srv.inShutdown.SetTrue()

	srv.mu.Lock()
	defer srv.mu.Unlock()

	lnerr := srv.listener.Close()
	srv.closeDoneChanLocked()

	ticker := time.NewTicker(shutdownPollInterval)
	defer ticker.Stop()
	for {
		if srv.getConnNum() == 0 {
			log.Info("server exit")
			return lnerr
		}
		select {
		case <-ctx.Done():
			log.Info("server exit")
			return ctx.Err()
		case <-ticker.C:
		}
	}
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

	srv.incrConnNum()

	return c
}

func (srv *Server) getConnNum() int64 {
	return atomic.LoadInt64(&srv.connNum)
}

func (srv *Server) incrConnNum() {
	atomic.AddInt64(&srv.connNum, 1)
}

func (srv *Server) decrConnNum() {
	atomic.AddInt64(&srv.connNum, -1)
}
