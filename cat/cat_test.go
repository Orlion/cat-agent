package cat

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/log"
)

var (
	hasInit bool
	initErr error
)

func testInit() error {
	if hasInit {
		return initErr
	}

	hasInit = true

	log.Init(&log.Config{
		StdoutLevel: "error",
	})

	initErr = Init(&config.Config{
		Domain:   "cat_agent_test_domain",
		Hostname: "cat_agent_test_hostname",
		Env:      "cat_agent_test_env",
		Ip:       "127.0.0.1",
		IpHex:    "",
		Servers:  []string{"127.0.0.1:8080"},
	})
	return initErr
}

func TestCreateMessageIdSingle(t *testing.T) {
	t.Log("TestCreateMessageIdSingle begin")

	if err := testInit(); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	messageId := CreateMessageId("cat-test-domain-single")
	if len(messageId) < 1 {
		t.Fatal("GetNextId get empty messageId")
	}

	t.Log("TestCreateMessageIdSingle end")
}

func TestCreateMessageIdThreeTimes(t *testing.T) {
	t.Log("TestCreateMessageIdThreeTimes begin")
	if err := testInit(); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	messageId1 := CreateMessageId("cat-test-domain-three_times")
	messageId2 := CreateMessageId("cat-test-domain-three_times")
	messageId3 := CreateMessageId("cat-test-domain-three_times")
	t.Logf("messageId1: %s, messageId2: %s, messageId3: %s", string(messageId1), string(messageId2), string(messageId3))
	if string(messageId1) == string(messageId2) {
		t.Fatalf("messageId1: %s == messageId2: %s", string(messageId1), string(messageId2))
	}

	if string(messageId1) == string(messageId3) {
		t.Fatalf("messageId1: %s == messageId3: %s", string(messageId1), string(messageId3))
	}

	if string(messageId2) == string(messageId3) {
		t.Fatalf("messageId2: %s == messageId3: %s", string(messageId2), string(messageId3))
	}

	t.Log("TestCreateMessageIdThreeTimes end")
}

func TestCreateMessageIdParallel(t *testing.T) {
	t.Log("TestCreateMessageIdParallel begin")

	if err := testInit(); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	type msgIdWithJ struct {
		msgId string
		j     int
	}

	msgIdCh := make(chan *msgIdWithJ, 50*10000)
	wg := &sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func(argsI, argsJ int) {
				messageId := CreateMessageId(fmt.Sprintf("cat-test-domain-%d", argsI))
				msgIdCh <- &msgIdWithJ{string(messageId), argsJ}
				wg.Done()
			}(i, j)
		}
	}

	wg.Wait()

	msgIdMap := make(map[string]int)

Loop:
	for {
		select {
		case msgIdWithJ := <-msgIdCh:
			if oldJ, exists := msgIdMap[msgIdWithJ.msgId]; exists {
				t.Fatalf("GetNextId get same messageId: %s, oldJ = %d, newJ = %d", msgIdWithJ.msgId, oldJ, msgIdWithJ.j)
			} else {
				msgIdMap[msgIdWithJ.msgId] = msgIdWithJ.j
			}
		default:
			break Loop
		}
	}
	t.Log("TestCreateMessageIdParallel end")
}
