package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ireuven89/hello-world/backend/environment"
)

func start() (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": environment.Variables.KafkaHost,
		"client.id":         "",
		"acks":              "all",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %v", err)
	}

	return producer, err
}
