package itemming

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/itemming/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

func TestMakeEndpointHealth(t *testing.T) {
	mockService := new(MockService)
	mockService.On("Health").Return(utils.ServiceHealthCheck{ServiceStatus: "OK"})

	endpoint := MakeEndpointHealth(mockService)
	response, err := endpoint(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, utils.ServiceHealthCheck{ServiceStatus: "OK"}, response)
}

func TestMakeEndpointGetItem(t *testing.T) {
	mockService := new(MockService)
	expectedItem := model.Item{ID: "test-uuid"}
	mockService.On("GetItem", "test-uuid").Return(expectedItem, nil)

	endpoint := MakeEndpointGetItem(mockService)
	response, err := endpoint(context.Background(), GetItemRequest{Uuid: "test-uuid"})

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, response)
}

func TestMakeEndpointListItems(t *testing.T) {
	mockService := new(MockService)
	expectedItems := []model.Item{{ID: "1"}, {ID: "2"}}
	mockService.On("GetItems", mock.Anything).Return(expectedItems, nil)

	endpoint := MakeEndpointListItems(mockService)
	response, err := endpoint(context.Background(), ListItemsRequest{})

	assert.NoError(t, err)
	assert.Equal(t, ListItemsResponse{items: expectedItems}, response)
}

func TestMakeEndpointCreateItem(t *testing.T) {
	mockService := new(MockService)
	mockService.On("CreateItem", mock.Anything).Return("new-uuid", nil)

	endpoint := MakeEndpointCreateItem(mockService)
	response, err := endpoint(context.Background(), CreateItemRequest{})

	assert.NoError(t, err)
	assert.Equal(t, CreateItemResponse{Uuid: "new-uuid"}, response)
}

func TestMakeEndpointCreateItems(t *testing.T) {
	mockService := new(MockService)
	mockService.On("CreateItems", mock.Anything).Return(nil)

	endpoint := MakeEndpointCreateItems(mockService)
	response, err := endpoint(context.Background(), CreateItemsRequest{})

	assert.NoError(t, err)
	assert.Nil(t, response)
}

func TestMakeEndpointUpdateItem(t *testing.T) {
	mockService := new(MockService)
	mockService.On("UpdateItem", mock.Anything).Return(nil)

	endpoint := MakeEndpointUpdateItem(mockService)
	response, err := endpoint(context.Background(), UpdateItemRequest{})

	assert.NoError(t, err)
	assert.Nil(t, response)
}

func TestMakeEndpointDeleteItem(t *testing.T) {
	mockService := new(MockService)
	mockService.On("DeleteItem", "test-uuid").Return(nil)

	endpoint := MakeEndpointDeleteItem(mockService)
	response, err := endpoint(context.Background(), GetItemRequest{Uuid: "test-uuid"})

	assert.NoError(t, err)
	assert.Nil(t, response)
}
