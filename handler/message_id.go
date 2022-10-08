package handler

import (
	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/server"
)

func CreateMessageId(req *server.Request) (status server.Status, payload []byte) {
	if len(req.Body) < 1 {
		status = server.StatusBadDomain
		return
	}

	payload = []byte(cat.CreateMessageId(string(req.Body)))

	return
}
