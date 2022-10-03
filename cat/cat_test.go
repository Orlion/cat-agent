package cat

import (
	"sync"
	"testing"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"go.uber.org/zap/zapcore"
)

var initErr error

func init() {
	log.Init(&log.Config{
		Level:    zapcore.DebugLevel,
		Filename: "",
	})

	initErr = Init(&config.Config{
		Domain:   "cat_test",
		Hostname: "cate_test_hostname",
		Env:      "cat_test_env",
		Ip:       "127.0.0.1",
		IpHex:    "",
		Servers:  []string{"127.0.0.1:8080"},
	})
}

func TestInit(t *testing.T) {
	t.Log("test Init begin")
	if initErr != nil {
		t.Fatalf("init error: %s", initErr)
	}
	t.Log("test init end")
}

func TestGetNextId(t *testing.T) {
	t.Log("test GetNextId begin")
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			messageId := GetNextId()
			t.Log(messageId)
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log("test GetNextId end")
}

func TestSend(t *testing.T) {
	t.Log("test Send begin")

	event := message.NewEvent("cat-test-event-type", "cat-test-event-name", message.SUCCESS, "cat-test-data", 100)

	tree := message.NewMessageTree()
	tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
	tree.SetThreadId(config.ThreadIdCatAgent)
	tree.SetThreadName(config.ThreadNameCatAgent)
	tree.SetMessageId(GetNextId())
	tree.SetParentMessageId("")
	tree.SetRootMessageId("")
	tree.SetMessage(event)
	tree.SetDiscard(false)
	Send(tree)
	Shutdown()
	t.Log("test Send end")
}
