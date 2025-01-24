package biddering

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/biddering/mocks"
	"github.com/ireuven89/hello-world/backend/biddering/model"
)

func TestBidderService_ListBidders(t *testing.T) {
	mockRepo := new(mocks.MockBidderRepo)
	logger := zap.NewNop()
	service := NewService(mockRepo, logger)

	input := model.BiddersInput{Item: "m"}
	expectedBidders := []model.Bidder{
		{Uuid: "1", Name: "Bidder1"},
		{Uuid: "2", Name: "Bidder2"},
	}

	mockRepo.On("List", input).Return(expectedBidders, nil)

	result, err := service.ListBidders(input)

	assert.NoError(t, err)
	assert.Equal(t, expectedBidders, result)
	mockRepo.AssertCalled(t, "List", input)
}

func TestBidderService_GetBidder(t *testing.T) {
	mockRepo := new(mocks.MockBidderRepo)
	logger := zap.NewNop()
	service := NewService(mockRepo, logger)

	uuid := "123"
	expectedBidder := model.Bidder{Uuid: "123", Name: "Test Bidder"}

	mockRepo.On("FindOne", uuid).Return(expectedBidder, nil)

	result, err := service.GetBidder(uuid)

	assert.NoError(t, err)
	assert.Equal(t, expectedBidder, result)
	mockRepo.AssertCalled(t, "FindOne", uuid)
}

func TestBidderService_CreateBidder(t *testing.T) {
	mockRepo := new(mocks.MockBidderRepo)
	logger := zap.NewNop()
	service := NewService(mockRepo, logger)

	input := model.BidderInput{Name: "New Bidder"}
	mockUUID := "mocks-uuid"

	mockRepo.On("Upsert", input).Return(mockUUID, nil)

	result, err := service.CreateBidder(input)

	assert.NoError(t, err)
	assert.Equal(t, mockUUID, result)
	mockRepo.AssertCalled(t, "Upsert", input)
}

func TestBidderService_UpdateBidder(t *testing.T) {
	mockRepo := new(mocks.MockBidderRepo)
	logger := zap.NewNop()
	service := NewService(mockRepo, logger)

	input := model.BidderInput{Uuid: "123", Name: "Updated Bidder"}
	mockUUID := "123"

	mockRepo.On("Upsert", input).Return(mockUUID, nil)

	result, err := service.UpdateBidder(input)

	assert.NoError(t, err)
	assert.Equal(t, mockUUID, result)
	mockRepo.AssertCalled(t, "Upsert", input)
}

func TestBidderService_Delete(t *testing.T) {
	mockRepo := new(mocks.MockBidderRepo)
	logger := zap.NewNop()
	service := NewService(mockRepo, logger)

	id := "123"

	mockRepo.On("Delete", id).Return(nil)

	err := service.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", id)
}
