package cat

import (
	"sync"
	"testing"

	"github.com/Orlion/cat-agent/cat/config"
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
			messageId := GetNextId("default")
			t.Log(messageId)
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log("test GetNextId end")
}
