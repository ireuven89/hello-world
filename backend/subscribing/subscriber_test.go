package subscribing

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	args2 := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return args2.Get(0).(chan amqp.Delivery), args2.Error(1)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	args2 := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return args2.Get(0).(amqp.Queue), args2.Error(1)
}

type MockAMQPConnection struct {
	mock.Mock
}

func (m *MockAMQPConnection) Channel() (*amqp.Channel, error) {
	args := m.Called()
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

/*
func TestSubscribe(t *testing.T) {
	// Mock the logger
	logger := zap.NewNop() // No-op logger for testing purposes

	// Mock the RabbitMQ connection and channel
	mockConn := new(MockAMQPConnection)
	mockCh := new(MockAMQPChannel)

	// Set up expectations for QueueDeclare
	mockConn.On("Channel").Return(mockCh, nil)
	mockCh.On("QueueDeclare", "test-queue", false, false, false, false, nil).Return(amqp.Queue{Name: "test-queue"}, nil)

	// Create a bidirectional channel to simulate message delivery
	deliveries := make(chan amqp.Delivery)

	// Mock the Consume method to return this channel
	mockCh.On("Consume", "test-queue", "", true, false, false, false, amqp.Table(nil)).Return(deliveries, nil)

	// Create the subscribing instance
	client := &Subscriber{
		conn:    mockConn,
		sub:     nil, // Assuming this is not relevant for the test
		logger:  logger,
		queue:   &amqp.Queue{Name: "test-queue"},
		channel: mockCh,
	}

	// Create the stop channel to control subscription termination
	stop := make(chan struct{})

	// Start the subscription in a goroutine
	done := make(chan struct{})
	go func() {
		err := client.Subscribe(stop) // This method will now run in a separate goroutine
		if err != nil {
			t.Logf("Subscription error: %v", err)
		}
		close(done) // Close the done channel when Subscribe completes
	}()

	// Simulate message delivery in a goroutine
	go func() {
		deliveries <- amqp.Delivery{Body: []byte("test message")}
		time.Sleep(1 * time.Second) // Simulate some delay before closing the channel
		close(deliveries)           // Close the channel to allow the Subscribe method to exit
	}()

	// Ensure that the subscribing finishes processing before moving forward
	select {
	case <-done: // Subscription completed successfully
		t.Log("Test completed successfully")
	case <-time.After(3 * time.Second): // Timeout after 3 seconds
		t.Fatal("Test timed out") // Fail the test if it takes too long
	}

	// Signal to stop the subscription gracefully
	close(stop)

	// Wait for the Subscribe method to finish gracefully
	<-done

	// Assert that expectations were met
	mockCh.AssertExpectations(t)
}

func TestSubscribeFailure(t *testing.T) {
	// Mock the logger
	logger := zap.NewNop()

	// Mock the RabbitMQ connection and channel
	mockConn := new(MockAMQPConnection)
	mockCh := new(MockAMQPChannel)

	deliveries := make(chan amqp.Delivery)
	// Set up expectations for QueueDeclare
	mockConn.On("Channel").Return(mockCh, nil)
	mockCh.On("QueueDeclare", "test-queue", false, false, false, false, nil).Return(amqp.Queue{Name: "test-queue"}, nil)

	// Simulate failure in Consume method
	mockCh.On("Consume", "test-queue", "", true, false, false, false, amqp.Table(nil)).Return(deliveries, amqp.ErrClosed)

	// Create the subscribing instance
	client := &Subscriber{
		conn:    mockConn,
		sub:     nil,
		logger:  logger,
		queue:   &amqp.Queue{Name: "test-queue"},
		channel: mockCh,
	}

	stop := make(chan struct{})

	// Call the Subscribe method
	err := client.Subscribe(stop)

	// Assertions
	assert.Error(t, err)
	mockCh.AssertExpectations(t) // Check if the Consume method was called
}*/
