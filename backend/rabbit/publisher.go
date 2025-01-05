package rabbit

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func (c *Client) Publish(message []byte) error {
	if err := c.channel.Publish("", c.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}); err != nil {
		c.logger.Error("failed to publish message: ", zap.Error(err))
		return err
	}

	return nil
}
