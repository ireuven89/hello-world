package rabbit

import (
	"github.com/brettallred/rabbit"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Client struct {
	conn    *amqp.Connection
	pub     *rabbit.Publisher
	sub     *rabbit.Subscriber
	logger  *zap.Logger
	queue   *amqp.Queue
	channel *amqp.Channel
}

func New(logger *zap.Logger) (*Client, error) {
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

	client := &Client{
		conn:    conn,
		pub:     &rabbit.Publisher{},
		sub:     &rabbit.Subscriber{Queue: queue.Name, Concurrency: 10},
		logger:  logger,
		queue:   &queue,
		channel: channel,
	}

	return client, nil
}
