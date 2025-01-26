package subscribing

import (
	"fmt"

	"github.com/brettallred/rabbit"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/environment"
)

type SService interface {
	Subscribe(stop chan struct{}) error
}

// AMQPConnection defines the interface for amqp.Connection
type AMQPConnection interface {
	Channel() (*amqp.Channel, error)
}

// AMQPChannel defines the interface for amqp.Channel
type AMQPChannel interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
}

type Subscriber struct {
	conn    AMQPConnection
	sub     *rabbit.Subscriber
	logger  *zap.Logger
	queue   *amqp.Queue
	channel AMQPChannel
}

func New(logger *zap.Logger) (SService, error) {
	url := environment.Variables.RabbitUrl
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(environment.Variables.RabbitQueue, false, false, false, false, nil)

	client := &Subscriber{
		conn:    conn,
		sub:     &rabbit.Subscriber{Queue: queue.Name},
		logger:  logger,
		queue:   &queue,
		channel: channel,
	}

	return client, nil
}

func (c *Subscriber) Subscribe(stop chan struct{}) error {

	messages, err := c.channel.Consume(c.queue.Name, "",
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to consume messages %v", err))
		return err
	}

	go func() {
		for d := range messages {
			c.logger.Info("Received message:", zap.Any("body", string(d.Body)))
		}
	}()

	// Listen for the stop signal or the forever condition
	// Wait for the stop signal to end the subscription
	select {
	case <-stop:
		// Gracefully stop if stop is received
		c.logger.Info("Stopped subscription gracefully.")
		return nil
	}

	c.logger.Info(fmt.Sprintf("started listening on queue %s", c.queue.Name))

	return nil
}
