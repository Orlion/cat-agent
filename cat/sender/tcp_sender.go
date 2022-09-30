package sender

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
)

type TcpSender struct {
	normal     chan *message.MessageTree
	high       chan *message.MessageTree
	config     *config.ConfigService
	wg         *sync.WaitGroup
	inShutdown atomicx.Bool
	running    bool
}

func NewTcpSender() *TcpSender {
	return &TcpSender{
		normal: make(chan *message.MessageTree, config.TcpSenderNormalQueueSize),
		high:   make(chan *message.MessageTree, config.TcpSenderHighQueueSize),
		config: config.GetInstance(),
		wg:     new(sync.WaitGroup),
	}
}

func (s *TcpSender) Run() {
	log.Info("tcp sender running...")

	for _, router := range s.config.GetRouters() {
		for i := 0; i < config.TcpSenderNormalQueueConsumerNum; i++ {
			s.wg.Add(1)
			go func(id int) {
				s.consume(router, id, 0)
				s.wg.Done()
			}(i)
		}

		for i := 0; i < config.TcpSenderHighQueueConsumerNum; i++ {
			s.wg.Add(1)
			go func(id int) {
				s.consume(router, id, 1)
				s.wg.Done()
			}(i)
		}
	}

	if !s.running {
		s.running = true
		go func() {
			for {
				// listen routers change
				s.config.RoutersCondWait()
				s.restart()
			}
		}()
	}
}

func (s *TcpSender) restart() {
	log.Info("tcp sender restart")
	s.inShutdown.SetTrue()
	s.wg.Wait()
	s.Run()
}

func (s *TcpSender) Shutdown() {
	log.Info("tcp sender shutdown...")

	s.inShutdown.SetTrue()

	close(s.normal)
	close(s.high)

	conn, err := net.DialTimeout("tcp", config.GetInstance().GetRouters()[0], time.Second)
	if err == nil {
		buf := make([]*message.MessageTree, config.TcpSenderQueueConsumerBufSize)

		for msg := range s.high {
			buf = append(buf, msg)
			if len(buf) == config.TcpSenderQueueConsumerBufSize {
				s.flush(conn, buf)
			}
		}

		for msg := range s.normal {
			buf = append(buf, msg)
			if len(buf) == config.TcpSenderQueueConsumerBufSize {
				s.flush(conn, buf)
			}
		}

		s.flush(conn, buf)

		conn.Close()
	}

	s.wg.Wait()
}

func (s *TcpSender) consume(server string, id int, chId int8) error {
	conn, err := net.DialTimeout("tcp", server, time.Second)
	if err != nil {
		log.Errorf("consumer try dial to %s error: %s", server, err.Error())
	}

	ticker := time.NewTicker(config.TcpSenderQueueConsumerTickerDuration)

	buf := make([]*message.MessageTree, 0, config.TcpSenderQueueConsumerBufSize)

	ch := s.normal
	if chId == 0 {
		log.Infof("tcp sender normal consumer: %d running...", id)
	} else {
		ch = s.high
		log.Infof("tcp sender high consumer: %d running...", id)
	}

	for !s.inShutdown.Get() {
		select {
		case msg := <-ch:
			buf = append(buf, msg)
			if len(buf) == config.TcpSenderQueueConsumerBufSize {
				if conn == nil {
					conn, err = net.DialTimeout("tcp", server, time.Second)
					if err != nil {
						log.Errorf("consumer retry dial to %s error: %s", server, err.Error())
					}
				}
				err = s.flush(conn, buf)
				if err != nil {
					log.Errorf("consumer flush %s error: %s", server, err.Error())
				}
				buf = buf[:0]
			}
		case <-ticker.C:
			s.flush(conn, buf)
			buf = buf[:0]
		}
	}

	s.flush(conn, buf)
	buf = nil

	conn.Close()

	ticker.Stop()

	return nil
}

func (s *TcpSender) flush(conn net.Conn, buf []*message.MessageTree) error {
	if len(buf) == 0 {
		return nil
	}
	// conn.SetWriteDeadline(time.Now().Add(time.Second))
	// todo
	fmt.Println(buf)
	// conn.Write(nil)
	return nil
}

func (s *TcpSender) Offer(tree *message.MessageTree) {
	if s.inShutdown.Get() {
		return
	}

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
