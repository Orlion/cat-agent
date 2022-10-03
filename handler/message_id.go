package handler

import (
	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/server"
)

func CreateMessageId(req *server.Request) (status server.Status, payload []byte) {
	payload = []byte(cat.GetNextId())
	return
}
