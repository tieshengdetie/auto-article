package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Consumer 定义消费者结构体
type Consumer struct {
	queueName    string
	exchangeName string
	routingKey   string
	handler      func([]byte)
	autoAck      bool
	consumerTag  string
	consumerType string // "simple", "topic", "fanout", "direct"
}

// ConsumerManagerInterface 定义消费者管理器接口
type ConsumerManagerInterface interface {
	AddSimpleConsumer(queueName string, handler func([]byte), autoAck bool) error
	AddTopicConsumer(exchangeName, queueName, routingKey string, handler func([]byte), autoAck bool) error
	AddFanoutConsumer(exchangeName, queueName string, handler func([]byte), autoAck bool) error
	AddDirectConsumer(exchangeName, queueName, routingKey string, handler func([]byte), autoAck bool) error
	StartConsumers(ctx context.Context) error
	StopConsumers() error
}

// ConsumerManager 实现消费者管理器接口
type ConsumerManager struct {
	client     *RabbitMQClient
	consumers  []*Consumer
	mu         sync.Mutex
	consumerWg sync.WaitGroup
	stopChan   chan struct{}
	stopped    bool
}

// NewConsumerManager 创建一个新的消费者管理器
func NewConsumerManager(client *RabbitMQClient) *ConsumerManager {
	return &ConsumerManager{
		client:    client,
		consumers: make([]*Consumer, 0),
		stopChan:  make(chan struct{}),
		stopped:   false,
	}
}

// AddSimpleConsumer 添加一个简单队列消费者
func (cm *ConsumerManager) AddSimpleConsumer(queueName, consumerTag string, handler func([]byte), autoAck bool) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return fmt.Errorf("consumer manager is stopped")
	}

	consumer := &Consumer{
		queueName:    queueName,
		handler:      handler,
		autoAck:      autoAck,
		consumerTag:  consumerTag,
		consumerType: TypeSimple,
	}
	cm.consumers = append(cm.consumers, consumer)
	return nil
}

// AddTopicConsumer 添加一个主题消费者
func (cm *ConsumerManager) AddTopicConsumer(exchangeName, queueName, consumerTag, routingKey string, handler func([]byte), autoAck bool) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return fmt.Errorf("consumer manager is stopped")
	}

	consumer := &Consumer{
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		handler:      handler,
		autoAck:      autoAck,
		consumerTag:  consumerTag,
		consumerType: TypeTopic,
	}
	cm.consumers = append(cm.consumers, consumer)
	return nil
}

// AddFanoutConsumer 添加一个扇出消费者
func (cm *ConsumerManager) AddFanoutConsumer(exchangeName, queueName, consumerTag string, handler func([]byte), autoAck bool) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return fmt.Errorf("consumer manager is stopped")
	}

	consumer := &Consumer{
		queueName:    queueName,
		exchangeName: exchangeName,
		handler:      handler,
		autoAck:      autoAck,
		consumerTag:  consumerTag,
		consumerType: TypeFanout,
	}
	cm.consumers = append(cm.consumers, consumer)
	return nil
}

// AddDirectConsumer 添加一个直连消费者
func (cm *ConsumerManager) AddDirectConsumer(exchangeName, queueName, consumerTag, routingKey string, handler func([]byte), autoAck bool) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return fmt.Errorf("consumer manager is stopped")
	}

	consumer := &Consumer{
		queueName:    queueName,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		handler:      handler,
		autoAck:      autoAck,
		consumerTag:  consumerTag,
		consumerType: TypeDirect,
	}
	cm.consumers = append(cm.consumers, consumer)
	return nil
}

// StartConsumers 启动所有消费者
func (cm *ConsumerManager) StartConsumers(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return fmt.Errorf("consumer manager is stopped")
	}

	for _, consumer := range cm.consumers {
		cm.consumerWg.Add(1)
		go func(c *Consumer) {
			defer cm.consumerWg.Done()
			var err error
			switch c.consumerType {
			case TypeSimple:
				err = cm.client.ConsumeFromQueue(c.queueName, c.consumerTag, c.handler, c.autoAck)
			case TypeTopic:
				err = cm.client.ConsumeWithTopic(c.exchangeName, c.queueName, c.consumerTag, c.routingKey, c.handler, c.autoAck)
			case TypeFanout:
				err = cm.client.ConsumeFanout(c.exchangeName, c.queueName, c.consumerTag, c.handler, c.autoAck)
			case TypeDirect:
				err = cm.client.ConsumeDirect(c.exchangeName, c.queueName, c.consumerTag, c.routingKey, c.handler, c.autoAck)
			default:
				log.Printf("Unknown consumer type: %s", c.consumerType)
				return
			}
			if err != nil {
				log.Printf("Failed to start consumer for %s: %s", c.consumerType, err)
			}
		}(consumer)
	}

	// 监听停止信号
	go func() {
		select {
		case <-ctx.Done():
			cm.StopConsumers()
		case <-cm.stopChan:
		}
	}()

	return nil
}

// StopConsumers 停止所有消费者
func (cm *ConsumerManager) StopConsumers() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.stopped {
		return nil
	}

	close(cm.stopChan)
	cm.consumerWg.Wait()
	cm.stopped = true
	return nil
}
