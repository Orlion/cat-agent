package config

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/atomicx"
	"github.com/Orlion/cat-agent/pkg/systemx"
)

type Config struct {
	Domain                       string   `yaml:"domain"`
	Hostname                     string   `yaml:"hostname"`
	Env                          string   `yaml:"env"`
	Ip                           string   `yaml:"ip"`
	IpHex                        string   `yaml:"ip_hex"`
	Servers                      []string `yaml:"servers"`
	SenderNormalQueueConsumerNum int      `yaml:"sender_normal_queue_consumer_num"`
	SenderHighQueueConsumerNum   int      `yaml:"sender_high_queue_consumer_num"`
}

type ConfigService struct {
	mu          sync.RWMutex
	config      *Config
	routers     []string
	routersCond *sync.Cond
	sample      float64
	enable      uint32
	done        chan struct{}
	wg          *sync.WaitGroup
}

func newConfigService(config *Config) (*ConfigService, error) {
	if err := withDefaultConf(config); err != nil {
		return nil, err
	}

	c := &ConfigService{
		config: config,
		done:   make(chan struct{}),
		wg:     new(sync.WaitGroup),
		enable: 1,
	}
	c.mu = sync.RWMutex{}
	c.routersCond = sync.NewCond(&c.mu)

	return c, nil
}

func (c *ConfigService) run() error {
	log.Info("config service running...")
	if err := c.pullRouters(); err != nil {
		return err
	}

	ticker := time.NewTicker(RouterUpdateDuration)

	c.wg.Add(1)
	go func() {
	Loop:
		for {
			select {
			case <-ticker.C:
				if err := c.pullRouters(); err != nil {
					log.Error(err.Error())
				}
			case <-c.done:
				break Loop
			}
		}
		c.wg.Done()
	}()

	return nil
}

func (c *ConfigService) shutdown() {
	log.Info("config service shutdown...")
	close(c.done)
	c.wg.Wait()
	log.Info("config service exit")
}

func (c *ConfigService) GetDomain() string {
	return c.config.Domain
}

func (c *ConfigService) GetHostname() string {
	return c.config.Hostname
}

func (c *ConfigService) GetEnv() string {
	return c.config.Env
}

func (c *ConfigService) GetIp() string {
	return c.config.Ip
}

func (c *ConfigService) GetIpHex() string {
	return c.config.IpHex
}

func (c *ConfigService) GetServers() []string {
	return c.config.Servers
}

func (c *ConfigService) GetSenderNormalQueueConsumerNum() int {
	return c.config.SenderNormalQueueConsumerNum
}

func (c *ConfigService) GetSenderHighQueueConsumerNum() int {
	return c.config.SenderHighQueueConsumerNum
}

func (c *ConfigService) GetRouters() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.routers
}

func (c *ConfigService) GetSample() float64 {
	return atomicx.LoadFloat64(&c.sample)
}

func (c *ConfigService) updateSample(v string) {
	sample, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Warnf("Sample should be a valid float, %s given", v)
	} else if math.Abs(sample-atomicx.LoadFloat64(&c.sample)) > 1e-9 {
		atomicx.StoreFloat64(&c.sample, sample)
		log.Infof("Sample rate has been set to %f%%", sample*100)
	}
}

func (c *ConfigService) updateRouters(router string) {
	newRouters := resolveServerAddresses(router)

	c.mu.Lock()
	defer c.mu.Unlock()
	oldLen, newLen := len(c.routers), len(newRouters)

	if newLen == 0 {
		log.Info("cannot established a connection to cat server")
		return
	} else if oldLen == 0 {
		log.Infof("routers has been initialized to: %v", newRouters)
		c.setRoutersLocked(newRouters)
	} else if oldLen != newLen {
		log.Infof("routers has been changed to: %s", newRouters)
		c.setRoutersLocked(newRouters)
	} else {
		for i := 0; i < oldLen; i++ {
			if c.routers[i] != newRouters[i] {
				log.Infof("routers has been changed to: %s", newRouters)
				c.setRoutersLocked(newRouters)
				break
			}
		}
	}
}

func (c *ConfigService) setRoutersLocked(routers []string) {
	c.routers = routers
	c.routersCond.Broadcast()
}

func (c *ConfigService) updateBlock(v string) {
	if v == "false" {
		if atomic.SwapUint32(&c.enable, 1) == 0 {
			log.Info("cat has been enabled")
		}
	} else {
		if atomic.SwapUint32(&c.enable, 0) == 1 {
			log.Info("cat has been disabled")
		}
	}
}

func (c *ConfigService) IsEnabled() bool {
	return atomic.LoadUint32(&c.enable) == 1
}

func (c *ConfigService) RoutersCondWait() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.routersCond.Wait()
}

var instance *ConfigService

func Init(config *Config) (err error) {
	instance, err = newConfigService(config)
	if err != nil {
		return
	}

	err = instance.run()
	if err != nil {
		return
	}

	return
}

func GetInstance() *ConfigService {
	return instance
}

func Shutdown() {
	instance.shutdown()
}

func withDefaultConf(config *Config) error {
	if config == nil {
		return errors.New("cat config cannot be empty")
	}

	if config.Domain == "" {
		return errors.New("domain cannot be empty")
	}

	var err error
	if config.Hostname, err = systemx.GetHostname(); err != nil {
		config.Hostname = DefaultHostname
	}

	if config.Env == "" {
		config.Env = DefaultEnv
	}

	if ip, err := systemx.GetLocalhostIp(); err != nil {
		log.Warnf("get localhost ip error: %s", err.Error())
		config.Ip = DefaultIp
		config.IpHex = DefaultIpHex
	} else {
		config.Ip = ip.String()
		config.IpHex = fmt.Sprintf("%02x%02x%02x%02x", ip[12], ip[13], ip[14], ip[15])
	}

	if len(config.Servers) < 1 {
		return errors.New("servers cannot be empty")
	}

	if config.SenderNormalQueueConsumerNum < 0 {
		return errors.New("sender normal queue consumer num cannot be less than 0")
	}

	if config.SenderNormalQueueConsumerNum == 0 {
		config.SenderNormalQueueConsumerNum = DefaultTcpSenderNormalQueueConsumerNum
	}

	if config.SenderHighQueueConsumerNum < 0 {
		return errors.New("sender high queue consumer num cannot be less than 0")
	}

	if config.SenderHighQueueConsumerNum == 0 {
		config.SenderHighQueueConsumerNum = DefaultTcpSenderHighQueueConsumerNum
	}

	return nil
}
