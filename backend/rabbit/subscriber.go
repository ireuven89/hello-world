package rabbit

import (
	"fmt"

	"go.uber.org/zap"
)

type Subscriber struct {
	logger *zap.Logger
}

func (c *Client) Subscribe() error {

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

	// Step 5: Handle messages
	forever := make(chan bool)

	go func() {
		for d := range messages {
			c.logger.Info("Received message:", zap.Any("body", string(d.Body)))
		}
	}()

	<-forever

	c.logger.Info(fmt.Sprintf("started listening on queue %s", c.queue.Name))

	return nil
}
