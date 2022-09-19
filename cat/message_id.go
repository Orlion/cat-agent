package cat

import (
	"fmt"
	"sync/atomic"
	"time"
)

func GetNextId(domain string) string {
	return messageIdFactory.getNextId(domain)
}

var messageIdFactory = new(MessageIdFactory)

type MessageIdFactory struct {
	index uint32
	hour  int
}

func (f *MessageIdFactory) getNextId(domain string) string {
	hour := int(time.Now().Unix() / 3600)
	if hour != f.hour {
		f.hour = hour
		currentIndex := atomic.LoadUint32(&f.index)
		atomic.CompareAndSwapUint32(&f.index, currentIndex, 0)
	}

	return fmt.Sprintf("%s-%s-%d-%d", domain, "todo:ipHex", hour, atomic.AddUint32(&f.index, 1))
}
