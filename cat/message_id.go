package cat

import (
	"bytes"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Orlion/cat-agent/cat/config"
)

type MessageIdFactory struct {
	mu        sync.RWMutex
	domain    string
	ipAddress []byte
	idPrefix  []byte
	index     uint32
	m         map[string]*uint32
	hour      []byte
}

func newMessageIdFactory() *MessageIdFactory {
	return &MessageIdFactory{
		domain:    config.GetInstance().GetDomain(),
		ipAddress: []byte(config.GetInstance().GetIpHex()),
		m:         make(map[string]*uint32),
	}
}

func (f *MessageIdFactory) getNextId(domain string) []byte {
	if domain == f.domain {
		return f.getLocalNextId()
	} else {
		return f.getDomainNextId(domain)
	}
}

func (f *MessageIdFactory) getLocalNextId() []byte {
	f.mu.RLock()
	defer f.mu.RUnlock()
	buf := new(bytes.Buffer)
	buf.Write(f.idPrefix)
	buf.WriteString(strconv.Itoa(int(atomic.AddUint32(&f.index, 1))))
	return buf.Bytes()
}

func (f *MessageIdFactory) getDomainNextId(domain string) []byte {
	buf := new(bytes.Buffer)

	f.mu.RLock()
	if index, exists := f.m[domain]; exists {
		buf.WriteString(domain)
		buf.WriteByte('-')
		buf.Write(f.ipAddress)
		buf.WriteByte('-')
		buf.Write(f.hour)
		buf.WriteByte('-')
		buf.WriteString(strconv.Itoa(int(atomic.AddUint32(index, 1))))
		f.mu.RUnlock()
	} else {
		f.mu.RUnlock()
		f.mu.Lock()
		if index, exists := f.m[domain]; exists {
			buf.WriteString(domain)
			buf.WriteByte('-')
			buf.Write(f.ipAddress)
			buf.WriteByte('-')
			buf.Write(f.hour)
			buf.WriteByte('-')
			buf.WriteString(strconv.Itoa(int(atomic.AddUint32(index, 1))))
		} else {
			buf.WriteString(domain)
			buf.WriteByte('-')
			buf.Write(f.ipAddress)
			buf.WriteByte('-')
			buf.Write(f.hour)
			buf.Write([]byte{'-', '1'})
			var i uint32 = 1
			f.m[domain] = &i
		}

		f.mu.Unlock()
	}

	return buf.Bytes()
}

func (f *MessageIdFactory) run() {

	ipAddressBuf := new(bytes.Buffer)
	ipAddressBuf.WriteString(f.domain)
	ipAddressBuf.WriteByte('-')
	ipAddressBuf.Write(f.ipAddress)
	ipAddressBuf.WriteByte('-')

	now := time.Now()

	f.mu.Lock()
	f.hour = []byte(strconv.FormatInt(now.Unix()/3600, 10))
	ipAddressBuf.Write(f.hour)
	ipAddressBuf.WriteByte('-')
	f.idPrefix = ipAddressBuf.Bytes()
	f.mu.Unlock()

	next := now.Add(time.Hour)
	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
	timer := time.NewTimer(next.Sub(now))

	go func() {
		for {
			<-timer.C

			ipAddressBuf := new(bytes.Buffer)
			ipAddressBuf.WriteString(f.domain)
			ipAddressBuf.WriteByte('-')
			ipAddressBuf.Write(f.ipAddress)
			ipAddressBuf.WriteByte('-')

			now = time.Now()

			f.mu.Lock()
			f.hour = []byte(strconv.FormatInt(now.Unix()/3600, 10))
			ipAddressBuf.Write(f.hour)
			ipAddressBuf.WriteByte('-')
			f.idPrefix = ipAddressBuf.Bytes()
			for _, index := range f.m {
				atomic.StoreUint32(index, 0)
			}
			f.mu.Unlock()

			next = now.Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
			timer.Reset(next.Sub(now))
		}
	}()
}
