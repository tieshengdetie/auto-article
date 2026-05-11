package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer
// @Description:
type Consumer struct {
	// consumer 用于消费消息
	readers map[string]*kafka.Reader // 使用 map 来保存 reader，key 是 topic
}

// NewConsumerWithSasl
//
//	@Description:
//	@param brokers
//	@param topic
//	@return *Consumer
func NewConsumerWithSasl(brokers []string, handlers map[string]RegisterHandlerStruct, groupIDs map[string]string, userName, password string) *Consumer {
	readers := make(map[string]*kafka.Reader)
	caCert, err := os.ReadFile("./only-4096-ca-cert")
	if err != nil {
		log.Fatal("读取 CA certificate 失败！ ", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("追加证书失败！")
	}
	// 创建TLS配置
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
	}
	// 创建 SASL PLAIN 认证机制
	var dialer kafka.Dialer
	dialer = kafka.Dialer{TLS: tlsConfig}
	dialer.SASLMechanism = plain.Mechanism{
		Username: userName,
		Password: password,
	}

	for key, registerHandler := range handlers {
		readers[key] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:                brokers,                 // broker地址 数组
			GroupID:                registerHandler.GroupId, // 消费者组id，每个消费者组可以消费kafka的完整数据，但是同一个消费者组中的消费者根据设置的分区消费策略共同消费kafka中的数据
			GroupTopics:            nil,
			Topic:                  registerHandler.Topic, // 消费哪个topic
			Partition:              0,
			Dialer:                 &dialer,
			QueueCapacity:          0,
			MinBytes:               0,
			MaxBytes:               0,
			MaxWait:                0,
			ReadBatchTimeout:       0,
			ReadLagInterval:        0,
			GroupBalancers:         nil,
			HeartbeatInterval:      0,
			CommitInterval:         time.Second, // offset 上报间隔
			PartitionWatchInterval: 0,
			WatchPartitionChanges:  false,
			SessionTimeout:         0,
			RebalanceTimeout:       0,
			JoinGroupBackoff:       0,
			RetentionTime:          0,
			StartOffset:            kafka.FirstOffset, // 仅对新创建的消费者组生效，从头开始消费，工作中可能更常用从最新的开始消费kafka.LastOffset
			ReadBackoffMin:         0,
			ReadBackoffMax:         0,
			Logger:                 nil,
			ErrorLogger:            nil,
			IsolationLevel:         0,
			MaxAttempts:            0,
			OffsetOutOfRangeError:  false,
		})

	}
	return &Consumer{readers: readers}
}

func NewConsumer(brokers []string, handlers map[string]RegisterHandlerStruct) *Consumer {
	readers := make(map[string]*kafka.Reader)
	for key, registerHandler := range handlers {

		readers[key] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:                brokers,                 // broker地址 数组
			GroupID:                registerHandler.GroupId, // 消费者组id，每个消费者组可以消费kafka的完整数据，但是同一个消费者组中的消费者根据设置的分区消费策略共同消费kafka中的数据
			GroupTopics:            nil,
			Topic:                  registerHandler.Topic, // 消费哪个topic
			Partition:              0,
			Dialer:                 nil,
			QueueCapacity:          0,
			MinBytes:               0,
			MaxBytes:               0,
			MaxWait:                0,
			ReadBatchTimeout:       0,
			ReadLagInterval:        0,
			GroupBalancers:         nil,
			HeartbeatInterval:      0,
			CommitInterval:         time.Second, // offset 上报间隔
			PartitionWatchInterval: 0,
			WatchPartitionChanges:  false,
			SessionTimeout:         0,
			RebalanceTimeout:       0,
			JoinGroupBackoff:       0,
			RetentionTime:          0,
			StartOffset:            kafka.LastOffset, // 仅对新创建的消费者组生效，从头开始消费，工作中可能更常用从最新的开始消费kafka.LastOffset
			ReadBackoffMin:         0,
			ReadBackoffMax:         0,
			Logger:                 nil,
			ErrorLogger:            nil,
			IsolationLevel:         0,
			MaxAttempts:            0,
			OffsetOutOfRangeError:  false,
		})

	}
	return &Consumer{readers: readers}
}

// StartConsumingWithHandlers ConsumeMessage 消费消息
func (k *Consumer) StartConsumingWithHandlers(ctx context.Context, registerHandlers map[string]RegisterHandlerStruct) {
	fmt.Println("等待读取Kafka消息")
	for key, reader := range k.readers {
		go func(k string, r *kafka.Reader) {
			for {
				msg, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("读取消息出错！: %v", err)
					continue
				}

				// 根据主题调用相应的消息处理器
				if registerHandler, ok := registerHandlers[k]; ok {
					registerHandler.Handler.HandleMessage(msg)
				} else {
					log.Printf("没有找到当前topic对应的处理方法！: %s", msg.Topic)
				}
			}
		}(key, reader)
	}

}

// Close 关闭Kafka工具类，释放资源
func (k *Consumer) Close() {
	for _, reader := range k.readers {
		reader.Close()
	}
}
