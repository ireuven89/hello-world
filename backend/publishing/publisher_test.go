package publishing

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/publishing/mocks"
)

func TestPublishMessage(t *testing.T) {
	// Mock the logger
	logger := zap.NewNop() // No-op logger for testing purposes

	// Mock the RabbitMQ connection and channel
	mockConn := new(mocks.MockConnection)
	mockCh := new(mocks.MockAMQPChannel)

	// Set up expectations
	mockConn.On("Channel").Return(mockCh, nil)
	mockCh.On("Publish", "", "test-queue", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("test message"),
	}).Return(nil)

	// Create the publishing instance
	client := &Publisher{
		conn:    mockConn,
		pub:     nil, // Assuming this is not relevant for the test
		logger:  logger,
		queue:   &amqp.Queue{Name: "test-queue"},
		channel: mockCh,
	}

	// Call the Publish method
	err := client.Publish([]byte("test message"))

	// Assertions
	assert.NoError(t, err)
	mockCh.AssertExpectations(t) // Check if the Publish method was called with correct arguments
}

func TestPublishMessageFailure(t *testing.T) {
	// Mock the logger
	logger := zap.NewNop()

	// Mock the RabbitMQ connection and channel
	mockConn := new(mocks.MockConnection)
	mockCh := new(mocks.MockAMQPChannel)

	// Set up expectations for failure
	mockConn.On("Channel").Return(mockCh, nil)
	mockCh.On("Publish", "", "test-queue", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("test message"),
	}).Return(amqp.ErrClosed)

	// Create the publishing instance
	client := &Publisher{
		conn:    mockConn,
		pub:     nil,
		logger:  logger,
		queue:   &amqp.Queue{Name: "test-queue"},
		channel: mockCh,
	}

	// Call the Publish method
	err := client.Publish([]byte("test message"))

	// Assertions
	assert.Error(t, err)
	mockCh.AssertExpectations(t) // Check if the Publish method was called correctly
}
