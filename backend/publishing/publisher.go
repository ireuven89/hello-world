package publishing

import (
	"github.com/brettallred/rabbit"
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
