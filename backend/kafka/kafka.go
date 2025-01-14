package kafka

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ireuven89/hello-world/backend/environment"
	"go.uber.org/zap"
)

type Service interface {
	Publish(string, interface{}) error
}

type Producer struct {
	logger   *zap.Logger
	producer *kafka.Producer
	kafka.TopicMetadata
}

func New() (*kafka.Producer, error) {

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

func (p *Producer) Publish(key string, message string) error {
	topic := "default-topic"
	m := kafka.Message{Key: []byte(key), Value: []byte(message), Timestamp: time.Now(), TopicPartition: kafka.TopicPartition{
		Topic:     &topic,
		Partition: kafka.PartitionAny,
	}}
	deliveryChan := make(chan kafka.Event)

	defer p.producer.Close()

	if err := p.producer.Produce(&m, deliveryChan); err != nil {
		return err
	}

	close(deliveryChan)
	return nil
}
