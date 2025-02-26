package producerring

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ireuven89/hello-world/backend/producerring/model"
)

// MigrationService handles migrations
type MigrationService struct {
	db         *mongo.Database
	kafkaTopic string
	producer   *kafka.Producer
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *mongo.Database, producer *kafka.Producer, topic string) *MigrationService {
	return &MigrationService{db: db, producer: producer, kafkaTopic: topic}
}

// CreateMigration initializes a new migration
func (s *MigrationService) CreateMigration(ctx context.Context, migration model.Migration) error {
	migration.Status = model.Pending.String()
	_, err := s.db.Collection("migrations").InsertOne(ctx, migration)
	return err
}

// PublishTask sends a task to Kafka
func (s *MigrationService) PublishTask(task model.MigrationTask) error {
	msg, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.kafkaTopic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
}
