package kafka

import (
	"AutoArticle/utils"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type dealTaskCenterMessage struct{}

const (
	taskCenterMessageKey = "taskCenterMessage"
)

func (h *dealTaskCenterMessage) HandleMessage(msg kafka.Message) {

	fmt.Printf("----开始消费---- \n")
	fmt.Printf("接收到的消息: %s \n", string(msg.Value))
	//Server.Broadcast(string(msg.Value))
	return
}

func init() {
	RegisterHandler(taskCenterMessageKey, "task_center_message", &dealTaskCenterMessage{}, generateGroupId("task_center_message"))
}

func generateGroupId(name string) (group string) {
	return name + "_" + utils.RandomString(5)
}
