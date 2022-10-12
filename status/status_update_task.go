package status

import (
	"bytes"
	"encoding/xml"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/cat/message"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/pkg/timex"
)

type StatusUpdateTask struct {
	statusExtensions []StatusExtension
}

func newStatusUpdateTask() *StatusUpdateTask {
	return &StatusUpdateTask{}
}

func (t *StatusUpdateTask) buildHeartbeat() {
	start := time.Now()
	trans := message.NewTransaction(config.TypeSystem, config.NameStatus, message.SUCCESS, "", timex.UnixMills(start), nil, 0)

	data, extensionTransList := t.buildExtension()
	for _, extensionTrans := range extensionTransList {
		trans.AddChild(extensionTrans)
	}

	heartbeat := message.NewHeartbeat(config.TypeHeartbeat, config.GetInstance().GetIp(), message.SUCCESS, data, timex.NowUnixMillis())
	trans.AddChild(heartbeat)
	trans.SetDurationInMicros(time.Now().Sub(start).Milliseconds())

	domain := config.GetInstance().GetDomain()
	tree := message.NewMessageTree()
	tree.SetMessage(trans)
	tree.SetDomain([]byte(domain))
	messageId := cat.CreateMessageId(domain)
	tree.SetMessageId(messageId)
	tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
	tree.SetThreadId([]byte(strconv.Itoa(os.Getpid())))
	tree.SetThreadName(config.ThreadNameCatAgent)
	tree.SetDiscard(false)

	cat.Send(tree)

	log.Infof("status update task send heartbeat, messageId: %s, ", messageId)
}

func (t *StatusUpdateTask) buildExtension() (string, []*message.Transaction) {
	extensionTransList := make([]*message.Transaction, 0, len(t.statusExtensions))
	status := Status{
		Extensions: make([]Extension, 0, len(t.statusExtensions)),
		CustomInfos: []CustomInfo{
			{"client", "cat-agent"},
			{"go-version", runtime.Version()},
		},
	}

	for _, statusExtension := range t.statusExtensions {
		start := time.Now()

		properties := statusExtension.GetProperties()
		if len(properties) > 0 {
			extension := Extension{
				Id:      statusExtension.GetId(),
				Desc:    statusExtension.GetDesc(),
				Details: make([]ExtensionDetail, 0),
			}

			for k, v := range statusExtension.GetProperties() {
				detail := ExtensionDetail{
					Id:    k,
					Value: v,
				}
				extension.Details = append(extension.Details, detail)
			}
			status.Extensions = append(status.Extensions, extension)
		}

		extensionTransList = append(extensionTransList, message.NewTransaction(config.TypeSystem, config.NameStatusExtensionPrefix+statusExtension.GetId(), message.SUCCESS, "", timex.UnixMills(start), nil, int64(time.Now().Sub(start).Milliseconds())))
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "\t")

	if err := encoder.Encode(status); err != nil {
		buf.Reset()
		buf.WriteString(err.Error())
		return buf.String(), extensionTransList
	}

	return buf.String(), extensionTransList
}

func (s *StatusUpdateTask) run() {
	log.Info("status update task running...")

	event := message.NewEvent(config.TypeSystem, config.NameReboot, message.SUCCESS, "", timex.NowUnixMillis())

	domain := config.GetInstance().GetDomain()
	tree := message.NewMessageTree()
	tree.SetMessage(event)
	tree.SetDomain([]byte(domain))
	messageId := cat.CreateMessageId(domain)
	tree.SetMessageId(messageId)
	tree.SetThreadGroupName(config.ThreadGroupNameCatAgent)
	tree.SetThreadId([]byte(strconv.Itoa(os.Getpid())))
	tree.SetThreadName(config.ThreadNameCatAgent)
	tree.SetDiscard(false)

	cat.Send(tree)
}

func (s *StatusUpdateTask) shutdown() {
	log.Info("status update task shutdown...")
	log.Info("status update task exit")
}

func Init() {

}
