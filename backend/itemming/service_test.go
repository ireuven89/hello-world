package itemming

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/itemming/mocks"
	"github.com/ireuven89/hello-world/backend/itemming/model"
)

func TestCreateItems(t *testing.T) {
	// Initialize the mock repository and logger
	mockRepo := new(mocks.MockRepo)
	mockLogger, _ := zap.NewDevelopment()

	service := &ServiceItem{
		repo:   mockRepo,
		logger: mockLogger,
	}

	// Prepare input data
	items := []model.ItemInput{{Uuid: "uuid"}, {Uuid: "item2"}} // Modify according to your model.ItemInput structure

	// Simulate failure first, then success
	mockRepo.On("BulkInsert", items).Return(errors.New("temporary failure")).Once() // First failure
	mockRepo.On("BulkInsert", items).Return(nil).Once()                             // Then success

	// Call the CreateItems method
	err := service.CreateItems(items)

	// Assert that the error is nil (should succeed after retries)
	assert.NoError(t, err)

	// Assert that the BulkInsert method was called twice
	mockRepo.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	// Initialize the mock repository and logger
	mockRepo := new(mocks.MockRepo)
	mockLogger, _ := zap.NewDevelopment()

	service := &ServiceItem{
		repo:   mockRepo,
		logger: mockLogger,
	}

	// Prepare input data
	uuid := "mock-uuid" // Modify according to your model.ItemInput structure

	// Simulate failure first, then success
	mockRepo.On("Delete", uuid).Return(errors.New("temporary failure")).Once() // First failure
	mockRepo.On("Delete", uuid).Return(nil).Once()                             // Then success

	// Call the CreateItems method
	err := service.DeleteItem(uuid)

	// Assert that the error is nil (should succeed after retries)
	assert.NoError(t, err)

	// Assert that the BulkInsert method was called twice
	mockRepo.AssertExpectations(t)

}

func TestServiceItem_CreateItem(t *testing.T) {
	// Initialize the mock repository and logger
	mockRepo := new(mocks.MockRepo)
	mockLogger, _ := zap.NewDevelopment()

	service := &ServiceItem{
		repo:   mockRepo,
		logger: mockLogger,
	}

	// Prepare input data
	invalidInput := model.ItemInput{} // Modify according to your model.ItemInput structure

	ValidInput := model.ItemInput{
		Uuid:        "uuid",
		Name:        "name",
		Description: "description",
		Price:       1,
	} // Modify according to your model.ItemInput structure

	// Simulate failure first, then success
	mockRepo.On("Insert", invalidInput).Return("", errors.New("failed inserting - missing mandatory fields")).Once() // fail
	mockRepo.On("Insert", ValidInput).Return("mock-uuid", nil).Once()                                                // Then success

	// Fail test
	id, err := service.CreateItem(invalidInput)
	assert.Error(t, err)
	assert.Equal(t, "", id)

	//Success test
	id, err = service.CreateItem(ValidInput)
	assert.NoError(t, err)
	assert.Equal(t, "mock-uuid", id)

	// Assert that the method was called twice
	mockRepo.AssertExpectations(t)

}
