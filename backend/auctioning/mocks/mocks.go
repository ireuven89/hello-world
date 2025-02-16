package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/auctioning/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type AuctionMockRepository struct {
	mock.Mock
}

type MockService struct {
	mock.Mock
}

func (mr *AuctionMockRepository) FindAll(req model.AuctionRequest) ([]model.Auction, error) {
	args := mr.Called(req)

	return args.Get(0).([]model.Auction), args.Error(1)
}
func (mr *AuctionMockRepository) FindOne(uuid string) (model.Auction, error) {
	args := mr.Called(uuid)

	return args.Get(0).(model.Auction), args.Error(1)
}
func (mr *AuctionMockRepository) Update(req model.AuctionRequest) error {
	args := mr.Called(req)

	return args.Error(0)
}
func (mr *AuctionMockRepository) Delete(uuid string) error {
	args := mr.Called(uuid)

	return args.Error(0)
}
func (mr *AuctionMockRepository) DbStatus() utils.DbStatus {
	args := mr.Called()

	return args.Get(0).(utils.DbStatus)
}

func (ms *MockService) Search(req model.AuctionRequest) ([]model.Auction, error) {
	args := ms.Called(req)

	return args.Get(0).([]model.Auction), args.Error(1)
}
func (ms *MockService) Find(uuid string) (model.Auction, error) {
	args := ms.Called(uuid)

	return args.Get(0).(model.Auction), args.Error(1)

}
func (ms *MockService) Delete(uuid string) error {
	args := ms.Called(uuid)

	return args.Error(0)
}
func (ms *MockService) Update(req model.AuctionRequest) error {
	args := ms.Called(req)

	return args.Error(0)
}
func (ms *MockService) Health() utils.ServiceHealthCheck {
	args := ms.Called()

	return args.Get(0).(utils.ServiceHealthCheck)
}
