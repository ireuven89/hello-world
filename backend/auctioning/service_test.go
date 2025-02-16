package auctioning

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/ireuven89/hello-world/backend/auctioning/mocks"
	"github.com/ireuven89/hello-world/backend/auctioning/model"
)

func TestSearch_Success(t *testing.T) {
	mockRepo := new(mocks.AuctionMockRepository) // Use existing mock
	logger := zaptest.NewLogger(t)
	auctionService := NewAuctionService(mockRepo, logger)

	request := model.AuctionRequest{Category: "Electronics"}
	expectedAuctions := []model.Auction{
		{Uuid: "1", Item: "Laptop"},
		{Uuid: "2", Item: "Phone"},
	}

	// Define mock expectations
	mockRepo.On("FindAll", request).Return(expectedAuctions, nil)

	// Call the method
	result, err := auctionService.Search(request)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedAuctions, result)
	mockRepo.AssertExpectations(t)
}

func TestSearch_Error(t *testing.T) {
	mockRepo := new(mocks.AuctionMockRepository)
	logger := zaptest.NewLogger(t)
	auctionService := NewAuctionService(mockRepo, logger)

	request := model.AuctionRequest{Category: "Furniture"}

	mockRepo.On("FindAll", request).Return([]model.Auction{}, errors.New("database error"))

	result, err := auctionService.Search(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestFindOne_Success(t *testing.T) {
	mockRepo := new(mocks.AuctionMockRepository) // Use existing mock
	testUUID := "123e4567-e89b-12d3-a456-426614174000"
	logger := zaptest.NewLogger(t)
	auctionService := NewAuctionService(mockRepo, logger)

	expectedAuction := model.Auction{
		ID:               testUUID,
		Item:             "Laptop",
		Price:            1000,
		WinningPrice:     1200,
		BiddersCount:     50,
		BiddersThreshold: 50,
	}

	// Define mock expectations
	mockRepo.On("FindOne", testUUID).Return(expectedAuction, nil)

	// Call the method
	result, err := auctionService.Find(testUUID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedAuction, result)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(mocks.AuctionMockRepository) // Mock repository
	mockLogger := zap.NewNop()                   // No-op logger for testing

	testUUID := "123e4567-e89b-12d3-a456-426614174000"

	// Mock the Delete function to return nil (success)
	mockRepo.On("Delete", testUUID).Return(nil)

	// Initialize the service with mock dependencies
	auctionService := NewAuctionService(mockRepo, mockLogger)

	// Execute the function
	err := auctionService.Delete(testUUID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
