package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/authenticating/model"
)

// MockService is a mocked implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Register(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockService) VerifyToken(jwtToken string) (string, error) {
	args := m.Called(jwtToken)
	return args.String(0), args.Error(1)
}

type InMemMock struct {
	Mock mock.Mock
}

func (ma *InMemMock) Save(username, password string) error {
	args := ma.Mock.Called(username, password)

	return args.Error(0)
}
func (ma *InMemMock) Find(username string) (model.User, error) {
	args := ma.Mock.Called(username)

	return args.Get(0).(model.User), args.Error(1)
}
