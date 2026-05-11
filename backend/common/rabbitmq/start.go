package rabbitmq

import (
	"AutoArticle/global"
	"fmt"
	"sync"
)

const (
	TypeSimple = "simple"
	TypeTopic  = "topic"
	TypeDirect = "direct"
	TypeFanout = "fanout"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(message []byte)
}

type MessageHandlerStruct struct {
	Handler      func([]byte)
	ExchangeName string
	QueueName    string
	ConsumerTag  string
	RoutingKey   string
	AutoAck      bool
	ConsumerType string
}

var MessageHandlerSyncLock sync.Mutex
var MessageHandlerSlice = make([]MessageHandlerStruct, 0)

func RegisterHandler(consumerType, exchangeName, queueName, consumerTag, routingKey string, handler MessageHandler, autoAck bool) {
	MessageHandlerSyncLock.Lock()
	defer MessageHandlerSyncLock.Unlock()
	MessageHandlerSlice = append(MessageHandlerSlice, MessageHandlerStruct{
		Handler:      handler.HandleMessage,
		ExchangeName: exchangeName,
		QueueName:    queueName,
		ConsumerTag:  consumerTag,
		RoutingKey:   routingKey,
		AutoAck:      autoAck,
		ConsumerType: consumerType,
	})

}
func StartRabbitMQConsumer() {

}

func PublishMessage(pushType, exchangeName, queueName, routingKey string, message []byte) (err error) {
	client := &RabbitMQClient{}
	config := global.ServerConfig.RabbitMq
	defer client.Close()
	mqUrl := fmt.Sprintf("amqp://%s:%s@%s:5672/readboyEbagMessageCenter", config.UserName, config.Password, config.Host)
	// 连接到 RabbitMQ 服务器
	err = client.Connect(mqUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}
	switch pushType {
	case TypeSimple:
		err = client.PublishSimple(queueName, message)
	case TypeTopic:
		err = client.PublishTopic(exchangeName, routingKey, message)
	case TypeDirect:
		err = client.PublishDirect(exchangeName, routingKey, message)
	case TypeFanout:
		err = client.PublishFanout(exchangeName, message)
	default:
		return fmt.Errorf("invalid consumer type: %s", pushType)
	}
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}

	return nil
}
