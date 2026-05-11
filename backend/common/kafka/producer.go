package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer
//
//	@Description: Kafka客户端
type Producer struct {
	// producer 用于生产消息
	producer *kafka.Writer
}

// ProducerClient
//
//	@Description:
//	@param brokers
//	@param topic
//	@return *Producer
func ProducerClient(brokers []string, topic string) *Producer {
	kafkaClient := &Producer{}
	// 创建生产者
	kafkaClient.producer = &kafka.Writer{
		Addr:                   kafka.TCP(brokers...), //TCP函数参数为不定长参数，可以传多个地址组成集群
		Topic:                  topic,
		Balancer:               &kafka.Hash{}, // 用于对key进行hash，决定消息发送到哪个分区
		MaxAttempts:            0,
		WriteBackoffMin:        0,
		WriteBackoffMax:        0,
		BatchSize:              0,
		BatchBytes:             0,
		BatchTimeout:           0,
		ReadTimeout:            0,
		WriteTimeout:           time.Second,       // kafka有时候可能负载很高，写不进去，那么超时后可以放弃写入，用于可以丢消息的场景
		RequiredAcks:           kafka.RequireNone, // 不需要任何节点确认就返回
		Async:                  false,
		Completion:             nil,
		Compression:            0,
		Logger:                 nil,
		ErrorLogger:            nil,
		Transport:              nil,
		AllowAutoTopicCreation: false, // 第一次发消息的时候，如果topic不存在，就自动创建topic，工作中禁止使用
	}

	return kafkaClient
}

// ProduceMessage 生产消息
func (k *Producer) ProduceMessage(message string) error {
	// 封装消息
	kafkaMessage := kafka.Message{
		Value: []byte(message),
	}
	// 发送消息
	return k.producer.WriteMessages(context.Background(), kafkaMessage)
}

// Close 关闭Kafka工具类，释放资源
func (k *Producer) Close() error {
	err := k.producer.Close()
	if err != nil {
		return err
	}
	return nil

}
