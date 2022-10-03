package handler

import (
	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/server"
)

func CreateMessageId(req *server.Request) (status server.Status, payload []byte) {
	if string(req.Body) != cat.GetDomain() {
		status = server.StatusBadDomain
		return
	}

	payload = []byte(cat.GetNextId())
	return
}
