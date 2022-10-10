package cat

import (
	"fmt"
	"regexp"
	"sync"
	"testing"

	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/log"
)

var (
	hasInit bool
)

func testInit(domain string) error {
	if hasInit {
		Shutdown()
	} else {
		log.Init(&log.Config{
			StdoutLevel: "debug",
		})
	}

	hasInit = true

	return Init(&config.Config{
		Domain:   domain,
		Hostname: "cat_agent_test_hostname",
		Env:      "cat_agent_test_env",
		Ip:       "127.0.0.1",
		IpHex:    "",
		Servers:  []string{"127.0.0.1:8080"},
	})
}

func TestCreateSingleLocalMessageId(t *testing.T) {
	t.Log("TestCreateSingleLocalMessageId begin")

	domain := "TestCreateSingleLocalMessageId"
	if err := testInit(domain); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	messageId := CreateMessageId(domain)
	t.Logf("messageId: %s", messageId)
	pattern := fmt.Sprintf(`%s-%s-\d+-\d+`, domain, config.GetInstance().GetIpHex())
	if !regexp.MustCompile(pattern).MatchString(string(messageId)) {
		t.Fatalf("messageId: %s not match: %s", messageId, pattern)
	}

	t.Log("TestCreateSingleLocalMessageId end")
}

func TestCreateSingleDomainMessageId(t *testing.T) {
	t.Log("TestCreateSingleDomainMessageId begin")

	if err := testInit("TestCreateSingleDomainMessageId"); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	domain := "test-domain"
	messageId := CreateMessageId(domain)
	t.Logf("messageId: %s", messageId)
	pattern := fmt.Sprintf(`%s-%s-\d+-\d+`, domain, config.GetInstance().GetIpHex())
	if !regexp.MustCompile(pattern).MatchString(string(messageId)) {
		t.Fatalf("messageId: %s not match: %s", messageId, pattern)
	}

	t.Log("TestCreateSingleDomainMessageId end")
}

func TestCreateThreeLocalMessageId(t *testing.T) {
	t.Log("TestCreateThreeLocalMessageId begin")
	domain := "TestCreateThreeLocalMessageId"
	if err := testInit(domain); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	messageId1 := CreateMessageId(domain)
	messageId2 := CreateMessageId(domain)
	messageId3 := CreateMessageId(domain)
	t.Logf("messageId1: %s, messageId2: %s, messageId3: %s", string(messageId1), string(messageId2), string(messageId3))
	pattern := fmt.Sprintf(`%s-%s-\d+-\d+`, domain, config.GetInstance().GetIpHex())
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(string(messageId1)) {
		t.Fatalf("messageId1: %s not match: %s", messageId1, pattern)
	}
	if !reg.MatchString(string(messageId2)) {
		t.Fatalf("messageId2: %s not match: %s", messageId2, pattern)
	}
	if !reg.MatchString(string(messageId3)) {
		t.Fatalf("messageId3: %s not match: %s", messageId3, pattern)
	}

	if string(messageId1) == string(messageId2) {
		t.Fatalf("messageId1: %s == messageId2: %s", string(messageId1), string(messageId2))
	}

	if string(messageId1) == string(messageId3) {
		t.Fatalf("messageId1: %s == messageId3: %s", string(messageId1), string(messageId3))
	}

	if string(messageId2) == string(messageId3) {
		t.Fatalf("messageId2: %s == messageId3: %s", string(messageId2), string(messageId3))
	}

	t.Log("TestCreateThreeLocalMessageId end")
}

func TestCreateThreeDomainMessageId(t *testing.T) {
	t.Log("TestCreateThreeDomainMessageId begin")
	if err := testInit("TestCreateThreeDomainMessageId"); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	domain := "test-domain"
	messageId1 := CreateMessageId(domain)
	messageId2 := CreateMessageId(domain)
	messageId3 := CreateMessageId(domain)
	t.Logf("messageId1: %s, messageId2: %s, messageId3: %s", string(messageId1), string(messageId2), string(messageId3))
	pattern := fmt.Sprintf(`%s-%s-\d+-\d+`, domain, config.GetInstance().GetIpHex())
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(string(messageId1)) {
		t.Fatalf("messageId1: %s not match: %s", messageId1, pattern)
	}
	if !reg.MatchString(string(messageId2)) {
		t.Fatalf("messageId2: %s not match: %s", messageId2, pattern)
	}
	if !reg.MatchString(string(messageId3)) {
		t.Fatalf("messageId3: %s not match: %s", messageId3, pattern)
	}

	if string(messageId1) == string(messageId2) {
		t.Fatalf("messageId1: %s == messageId2: %s", string(messageId1), string(messageId2))
	}

	if string(messageId1) == string(messageId3) {
		t.Fatalf("messageId1: %s == messageId3: %s", string(messageId1), string(messageId3))
	}

	if string(messageId2) == string(messageId3) {
		t.Fatalf("messageId2: %s == messageId3: %s", string(messageId2), string(messageId3))
	}

	t.Log("TestCreateThreeDomainMessageId end")
}

func TestCreateLocalMessageIdParallel(t *testing.T) {
	t.Log("TestCreateLocalMessageIdParallel begin")

	baseDomain := "TestCreateLocalMessageIdParallel-"
	if err := testInit(baseDomain + "0"); err != nil {
		t.Fatalf("testInit error: %s", err)
	}

	type msgId struct {
		msgId string
		i     int
		j     int
	}

	msgIdCh := make(chan *msgId, 50*10000)
	wg := &sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		for j := 0; j < 10000; j++ {
			wg.Add(1)
			go func(argsI, argsJ int) {
				messageId := CreateMessageId(fmt.Sprintf("%s-%d", baseDomain, argsI))
				msgIdCh <- &msgId{string(messageId), argsI, argsJ}
				wg.Done()
			}(i, j)
		}
	}

	wg.Wait()

	msgIdMap := make(map[string]int)

Loop:
	for {
		select {
		case msgId := <-msgIdCh:
			pattern := fmt.Sprintf(`%s-%d-%s-\d+-\d+`, baseDomain, msgId.i, config.GetInstance().GetIpHex())
			if !regexp.MustCompile(pattern).MatchString(string(msgId.msgId)) {
				t.Fatalf("messageId1: %s not match: %s", msgId.msgId, pattern)
			}

			if oldJ, exists := msgIdMap[msgId.msgId]; exists {
				t.Fatalf("GetNextId get same messageId: %s, oldJ = %d, newJ = %d", msgId.msgId, oldJ, msgId.j)
			} else {
				msgIdMap[msgId.msgId] = msgId.j
			}
		default:
			break Loop
		}
	}
	t.Log("TestCreateMessageIdParallel end")
}
