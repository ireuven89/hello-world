package publishing

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	args2 := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return args2.Get(0).(amqp.Queue), args2.Error(1)
}

type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) Channel() (*amqp.Channel, error) {
	args := m.Called()
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

func TestPublishMessage(t *testing.T) {
	// Mock the logger
	logger := zap.NewNop() // No-op logger for testing purposes

	// Mock the RabbitMQ connection and channel
	mockConn := new(MockConnection)
	mockCh := new(MockAMQPChannel)

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
	mockConn := new(MockConnection)
	mockCh := new(MockAMQPChannel)

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
