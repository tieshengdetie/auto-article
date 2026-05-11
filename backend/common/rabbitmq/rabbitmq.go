package rabbitmq

import (
	"github.com/streadway/amqp"
)

// RabbitMQClient 定义 RabbitMQ 客户端
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect 连接到 RabbitMQ 服务器
func (c *RabbitMQClient) Connect(url string) error {
	var err error
	c.conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭 RabbitMQ 连接
func (c *RabbitMQClient) Close() error {
	if c.channel != nil {
		err := c.channel.Close()
		if err != nil {
			return err
		}
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// ConsumeFromQueue 从队列消费消息（简单队列模式）
func (c *RabbitMQClient) ConsumeFromQueue(queueName, consumerTag string, handler func([]byte), autoAck bool) error {
	queue, err := c.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

// ConsumeWithTopic 使用主题模式消费消息
func (c *RabbitMQClient) ConsumeWithTopic(exchangeName, queueName, consumerTag, routingKey string, handler func([]byte), autoAck bool) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		queue.Name,   // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

// ConsumeFanout 使用扇出模式消费消息
func (c *RabbitMQClient) ConsumeFanout(exchangeName, queueName, consumerTag string, handler func([]byte), autoAck bool) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

// ConsumeDirect 使用直连模式消费消息
func (c *RabbitMQClient) ConsumeDirect(exchangeName, queueName, consumerTag, routingKey string, handler func([]byte), autoAck bool) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		queue.Name,   // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		autoAck,     // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

// PublishSimple 发布消息到简单队列
func (c *RabbitMQClient) PublishSimple(queueName string, body []byte) error {
	queue, err := c.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

// PublishTopic 发布消息到主题模式
func (c *RabbitMQClient) PublishTopic(exchangeName, routingKey string, body []byte) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

// PublishFanout 发布消息到扇出模式
func (c *RabbitMQClient) PublishFanout(exchangeName string, body []byte) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

// PublishDirect 发布消息到直连模式
func (c *RabbitMQClient) PublishDirect(exchangeName, routingKey string, body []byte) error {
	err := c.channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
