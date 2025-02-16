package mocks

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
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
