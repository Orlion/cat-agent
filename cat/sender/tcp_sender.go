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
				newConsumer(id, router, "normal", s.normal).run()
				s.wg.Done()
			}(i)
		}

		for i := 0; i < config.TcpSenderHighQueueConsumerNum; i++ {
			s.wg.Add(1)
			go func(id int) {
				newConsumer(id, router, "high", s.high).run()
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
	s.Shutdown()
	s.Run()
}

func (s *TcpSender) Shutdown() {
	log.Info("tcp sender shutdown...")

	s.inShutdown.SetTrue()

	close(s.normal)
	close(s.high)

	s.wg.Wait()

	log.Info("tcp sender exit")
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

type Consumer struct {
	name     string
	server   string
	ch       <-chan *message.MessageTree
	conn     net.Conn
	connTime time.Time
	buf      []*message.MessageTree
}

func newConsumer(id int, server, chName string, ch <-chan *message.MessageTree) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("%s-%s-%d", chName, server, id),
		server: server,
		ch:     ch,
		buf:    make([]*message.MessageTree, 0, config.TcpSenderQueueConsumerBufSize),
	}
}

func (c *Consumer) run() {
	log.Infof("consumer %s running...", c.name)

	ticker := time.NewTicker(config.TcpSenderQueueConsumerTickerDuration)

Loop:
	for {
		select {
		case msg, ok := <-c.ch:
			if !ok {
				break Loop
			}
			c.buf = append(c.buf, msg)
			if len(c.buf) == config.TcpSenderQueueConsumerBufSize {
				c.flush(false)
			}
		case <-ticker.C:
			c.flush(false)
		}
	}

	ticker.Stop()

	c.flush(true)
	c.buf = nil

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) connect(nonblock bool) error {
	var (
		err       error
		tempDelay time.Duration
	)

	if c.conn != nil && time.Now().Sub(c.connTime) < 10*time.Minute {
		return nil
	}

	c.conn = nil

	for {
		c.conn, err = net.DialTimeout("tcp", c.server, time.Second)
		if err == nil {
			c.connTime = time.Now()
			break
		}

		if nonblock {
			return err
		}

		log.Errorf("consumer %s dial to %s error: %s, tempDelay: %d", c.name, c.server, err, tempDelay)

		if tempDelay == 0 {
			tempDelay = 100 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 5 * time.Second; tempDelay > max {
			tempDelay = max
		}

		time.Sleep(tempDelay)
	}

	return nil
}

func (c *Consumer) flush(nonblock bool) {
	if len(c.buf) == 0 {
		return
	}
	if err := c.connect(nonblock); err != nil {
		return
	}
	// conn.SetWriteDeadline(time.Now().Add(time.Second))
	// todo
	fmt.Println(c.buf)
	c.buf = c.buf[:0]
}
