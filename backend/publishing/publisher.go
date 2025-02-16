package publishing

import (
	"context"
	"fmt"
	"time"

	"github.com/brettallred/rabbit"
	"github.com/sethvargo/go-retry"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/environment"
)

type PService interface {
	Publish(message []byte) error
}

type AMQPConnection interface {
	Channel() (*amqp.Channel, error)
}

// AMQPChannel defines the interface for amqp.Channel
type AMQPChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
}

type Publisher struct {
	conn    AMQPConnection
	pub     *rabbit.Publisher
	logger  *zap.Logger
	queue   *amqp.Queue
	channel AMQPChannel
}

func New(logger *zap.Logger) (PService, error) {
	var queue amqp.Queue
	var conn AMQPConnection
	var channel AMQPChannel
	retryFunc := retry.NewConstant(time.Second * 2)
	ret := retry.WithMaxRetries(4, retryFunc)

	err := retry.Do(context.Background(), ret, func(ctx context.Context) error {
		url := environment.Variables.RabbitUrl
		conn, err := amqp.Dial(url)

		if err != nil {
			logger.Error(fmt.Sprintf("failed dialing %v", err))
			return err
		}

		channel, err := conn.Channel()

		if err != nil {
			logger.Error("failed creating channel")
			return err
		}

		queue, err = channel.QueueDeclare(environment.Variables.RabbitQueue, false, false, false, false, nil)

		if err != nil {
			logger.Error("failed queue declare")
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("failed connection to rabbit")
		return nil, err
	}
	client := &Publisher{
		conn:    conn,
		pub:     &rabbit.Publisher{},
		logger:  logger,
		queue:   &queue,
		channel: channel,
	}

	return client, nil
}

func (p *Publisher) Publish(message []byte) error {
	if err := p.channel.Publish("", p.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}); err != nil {
		p.logger.Error("failed to publish message: ", zap.Error(err))
		return err
	}

	return nil
}
