package sender

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
)

type TcpSender struct {
	normal  chan *message.MessageTree
	high    chan *message.MessageTree
	config  *config.ConfigService
	running bool
	wg      *sync.WaitGroup
	cancel  func()
}

func NewTcpSender() *TcpSender {
	return &TcpSender{
		normal: make(chan *message.MessageTree, config.NormalPriorityQueueSize),
		high:   make(chan *message.MessageTree, config.HighPriorityQueueSize),
		config: config.GetInstance(),
		wg:     new(sync.WaitGroup),
	}
}

func (s *TcpSender) Run() {
	var ctx context.Context

	ctx, s.cancel = context.WithCancel(context.Background())
	for _, router := range s.config.GetRouters() {
		for i := 0; i < config.QueueConsumerNum; i++ {
			s.wg.Add(1)
			go func(id int) {
				s.consume(ctx, router, id, s.normal)
				s.wg.Done()
			}(i)
		}
	}

	if !s.running {
		s.running = true
		go func() {
			// listen routers change
			s.config.RoutersCondWait()
			s.restart()
		}()
	}
}

func (s *TcpSender) restart() {
	s.Shutdown()
	s.Run()
}

func (s *TcpSender) Shutdown() {
	s.cancel()
	s.wg.Wait()
}

func (s *TcpSender) consume(ctx context.Context, server string, id int, ch chan *message.MessageTree) error {
	conn, err := net.DialTimeout("tcp", server, time.Second)
	if err != nil {
		return err
	}

	timer := time.NewTimer(config.QueueConsumerTimerDuration)

	buf := make([]*message.MessageTree, config.QueueConsumerBufSize)

	log.Infof("tcp sender consumer: %d running...", id)

Loop:
	for {
		select {
		case msg := <-ch:
			buf = append(buf, msg)
			if len(buf) == config.QueueConsumerBufSize {
				s.flush(conn, buf)
				buf = buf[:0]
			}
		case <-timer.C:
			if len(buf) > 0 {
				s.flush(conn, buf)
				buf = buf[:0]
			}
		case <-ctx.Done():
			if len(buf) > 0 {
				s.flush(conn, buf)
			}
			break Loop
		}
	}

	return nil
}

func (s *TcpSender) flush(conn net.Conn, buf []*message.MessageTree) {
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	// todo
	fmt.Println(buf)
	conn.Write(nil)
}

func (s *TcpSender) Offer(tree *message.MessageTree) {
	if tree.GetMessage().IsSuccess() {
		select {
		case s.normal <- tree:
		default:

		}
	} else {
		select {
		case s.high <- tree:
		default:

		}
	}
}
