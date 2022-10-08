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
	mu     sync.RWMutex
	indexs map[string]*uint32
	hour   int
}

func newMessageIdFactory() *MessageIdFactory {
	return &MessageIdFactory{
		indexs: make(map[string]*uint32),
	}
}

func (f *MessageIdFactory) getNextId(domain string) []byte {
	buf := new(bytes.Buffer)

	f.mu.RLock()
	if index, exists := f.indexs[domain]; exists {
		buf.WriteString(domain)
		buf.WriteByte('-')
		buf.WriteString(config.GetInstance().GetIpHex())
		buf.WriteByte('-')
		buf.WriteString(strconv.Itoa(f.hour))
		buf.WriteByte('-')
		buf.WriteString(strconv.Itoa(int(atomic.AddUint32(index, 1))))
		f.mu.RUnlock()
	} else {
		f.mu.RUnlock()
		f.mu.Lock()
		if index, exists := f.indexs[domain]; exists {
			buf.WriteString(domain)
			buf.WriteByte('-')
			buf.WriteString(config.GetInstance().GetIpHex())
			buf.WriteByte('-')
			buf.WriteString(strconv.Itoa(f.hour))
			buf.WriteByte('-')
			buf.WriteString(strconv.Itoa(int(atomic.AddUint32(index, 1))))
		} else {
			buf.WriteString(domain)
			buf.WriteByte('-')
			buf.WriteString(config.GetInstance().GetIpHex())
			buf.WriteByte('-')
			buf.WriteString(strconv.Itoa(f.hour))
			buf.Write([]byte{'-', '1'})
			var i uint32 = 1
			f.indexs[domain] = &i
		}

		f.mu.Unlock()
	}

	return buf.Bytes()
}

func (f *MessageIdFactory) run() {
	now := time.Now()

	f.mu.Lock()
	f.hour = int(now.Unix()) / 3600
	f.mu.Unlock()

	next := now.Add(time.Hour)
	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
	timer := time.NewTimer(next.Sub(now))
	go func() {
		for {
			<-timer.C

			now = time.Now()
			next = now.Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
			timer.Reset(next.Sub(now))

			f.mu.Lock()
			f.hour = int(now.Unix()) / 3600
			for _, index := range f.indexs {
				atomic.StoreUint32(index, 0)
			}
			f.mu.Unlock()
		}
	}()
}
