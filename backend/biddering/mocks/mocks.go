package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/biddering/model"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) ListBidders(input model.BiddersInput) ([]model.Bidder, error) {
	args := m.Called(input)
	return args.Get(0).([]model.Bidder), args.Error(1)
}

func (m *MockService) GetBidder(uuid string) (model.Bidder, error) {
	args := m.Called(uuid)
	return args.Get(0).(model.Bidder), args.Error(1)
}

func (m *MockService) CreateBidder(input model.BidderInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockService) UpdateBidder(input model.BidderInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockService) DeleteBidder(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockRedis struct {
	Mock mock.Mock
}

func (mdb *MockRedis) Set(key string, value interface{}, ttl time.Duration) error {
	args := mdb.Mock.Called(key, value, ttl)

	return args.Error(0)
}

func (mdb *MockRedis) Get(key string) (interface{}, error) {
	args := mdb.Mock.Called(key)

	return args.Get(0), args.Error(1)
}

type MockBidderRepo struct {
	mock.Mock
}

func (m *MockBidderRepo) List(input model.BiddersInput) ([]model.Bidder, error) {
	args := m.Called(input)
	return args.Get(0).([]model.Bidder), args.Error(1)
}

func (m *MockBidderRepo) FindOne(uuid string) (model.Bidder, error) {
	args := m.Called(uuid)
	return args.Get(0).(model.Bidder), args.Error(1)
}

func (m *MockBidderRepo) Upsert(input model.BidderInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockBidderRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
