package handler

import (
	"github.com/Orlion/cat-agent/server"
)

func SendMessage(req *server.Request) (status server.Status, payload []byte) {
	// read header
	r := &messageReader{
		len:  len(req.Body),
		body: req.Body,
	}

	_, err := r.readHeader()
	if err != nil {

	}

	return
}

type messageHeader struct {
	domain          []byte
	threadGroupName []byte
	threadId        []byte
	threadName      []byte
	messageId       []byte
	parentMessageId []byte
	rootMessageId   []byte
}

type messageReader struct {
	i    int
	len  int
	body []byte
}

func (r *messageReader) readHeader() (header *messageHeader, err error) {
	var (
		end bool
		eof bool
	)

	header = new(messageHeader)

	header.domain, end, eof = r.readElement()
	if end {
		return
	}
	if eof {
		return
	}

	header.threadGroupName, end, eof = r.readElement()
	if end {
		return
	}

	header.threadId, end, eof = r.readElement()
	if end {
		return
	}

	header.threadName, end, eof = r.readElement()
	if end {
		return
	}

	header.messageId, end, eof = r.readElement()
	if end {
		return
	}

	header.parentMessageId, end, eof = r.readElement()
	if end {
		return
	}

	header.rootMessageId, end, eof = r.readElement()

	return
}

func (r *messageReader) readElement() (b []byte, end bool, eof bool) {
	b = make([]byte, 0)

	for {
		if r.i == r.len-1 {
			eof = true
			break
		}

		if r.body[r.i] == '\n' {
			end = true
			break
		}
		if r.body[r.i] == '\t' {
			break
		}
		b = append(b, r.body[r.i])
		r.i++
	}

	return
}
