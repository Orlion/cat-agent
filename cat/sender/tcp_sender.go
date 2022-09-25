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
	normal chan *message.MessageTree
	high   chan *message.MessageTree
	config *config.ConfigService
}

func NewTcpSender() *TcpSender {
	return &TcpSender{
		normal: make(chan *message.MessageTree, config.NormalPriorityQueueSize),
		high:   make(chan *message.MessageTree, config.HighPriorityQueueSize),
		config: config.GetInstance(),
	}
}

func (s *TcpSender) Run() {
	var (
		wg = new(sync.WaitGroup)
	)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		for _, router := range s.config.GetRouters() {
			for i := 0; i < config.QueueConsumerNum; i++ {
				wg.Add(1)
				go func() {
					s.consume(ctx, router, i, s.normal)
					wg.Done()
				}()
			}
		}

		// wait for routers change
		s.config.RoutersCondWait()
		// if routers changed cancel consume
		cancel()
		// wait all the consumer stop
		wg.Wait()
	}
}

func (s *TcpSender) consume(ctx context.Context, server string, i int, ch chan *message.MessageTree) error {
	conn, err := net.DialTimeout("tcp", server, time.Second)
	if err != nil {
		return err
	}

	timer := time.NewTimer(config.QueueConsumerTimerDuration)

	buf := make([]*message.MessageTree, config.QueueConsumerBufSize)

	log.Infof("tcp sender consumer: %d running...", i)

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
