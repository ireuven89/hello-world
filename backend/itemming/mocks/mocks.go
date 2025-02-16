package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/itemming/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockRedis) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value)

	return args.Error(0)
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) BulkInsert(items []model.ItemInput) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *MockRepo) ListItems(input model.ListInput) ([]model.Item, error) {
	args := m.Called(input)

	return args.Get(0).([]model.Item), args.Error(1)
}
func (m *MockRepo) GetItem(uuid string) (model.Item, error) {
	args := m.Called(uuid)

	return args.Get(0).(model.Item), args.Error(1)
}
func (m *MockRepo) Upsert(input model.ItemInput) (string, error) {
	args := m.Called(input)

	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepo) Delete(uuid string) error {
	args := m.Called(uuid)

	return args.Error(0)
}

func (m *MockRepo) Insert(input model.ItemInput) (string, error) {
	args := m.Called(input)

	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepo) Update(input model.ItemInput) error {
	args := m.Called(input)

	return args.Error(0)
}

func (m *MockRepo) DBstatus() utils.DbStatus {

	return utils.DbStatus{}
}
