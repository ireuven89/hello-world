package mocks

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
