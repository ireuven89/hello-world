package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/userring/model"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListUsers(input model.UserFetchInput) ([]model.User, error) {
	args := m.Called(input)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockRepository) FindUser(uuid string) (model.User, error) {
	args := m.Called(uuid)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockRepository) Upsert(input model.UserUpsertInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) Delete(uuid string) error {
	args := m.Called(uuid)
	return args.Error(0)
}

type MockService struct {
	mock.Mock
}

func (ms *MockService) ListUsers(input model.UserFetchInput) ([]model.User, error) {
	args := ms.Called(input)

	return args.Get(0).([]model.User), args.Error(1)
}
func (ms *MockService) GetUser(uuid string) (model.User, error) {
	args := ms.Called(uuid)

	return args.Get(0).(model.User), args.Error(1)
}
func (ms *MockService) CreateUser(input model.UserUpsertInput) (string, error) {
	args := ms.Called(input)

	return args.Get(0).(string), args.Error(1)
}
func (ms *MockService) UpdateUser(input model.UserUpsertInput) error {
	args := ms.Called(input)

	return args.Error(0)
}
func (ms *MockService) DeleteUser(uuid string) error {
	args := ms.Called(uuid)

	return args.Error(0)
}
