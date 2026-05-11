package kafka

import (
	"AutoArticle/global"
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(message kafka.Message)
}
type HandlerMap map[string]RegisterHandlerStruct

var GroupIdMap = map[string]string{}

var TopicHandlerMapping = make(HandlerMap)

type RegisterHandlerStruct struct {
	Topic   string
	Handler MessageHandler
	GroupId string
}

func RegisterHandler(key string, topic string, handler MessageHandler, groupId string) {
	if _, ok := TopicHandlerMapping[key]; !ok {
		TopicHandlerMapping[key] = RegisterHandlerStruct{Topic: topic, Handler: handler, GroupId: groupId}
	} else {
		panic("handler key already registered")
	}

}
func StartConsumer() *Consumer {
	config := global.ServerConfig.Kafka
	brokers := strings.Split(config.Brokers, ",")
	// 创建 Kafka 消费者：根据 handlers 的 key 和对应的 groupID 创建消费者
	consumer := NewConsumer(brokers, TopicHandlerMapping)
	// 启动 Kafka 消费服务
	fmt.Println("启动 Kafka 消费服务！")
	go func() {
		ctx := context.Background()
		consumer.StartConsumingWithHandlers(ctx, TopicHandlerMapping)
	}()

	return consumer
}
func StartConsumerWithSasl() *Consumer {
	config := global.ServerConfig.Kafka
	brokers := strings.Split(",", config.Brokers)
	userName := config.User
	password := config.Password
	// 创建 Kafka 消费者：根据 handlers 的 key 和对应的 groupID 创建消费者
	consumer := NewConsumerWithSasl(brokers, TopicHandlerMapping, GroupIdMap, userName, password)
	// 启动 Kafka 消费服务
	go func() {
		ctx := context.Background()
		consumer.StartConsumingWithHandlers(ctx, TopicHandlerMapping)
	}()

	return consumer
}
