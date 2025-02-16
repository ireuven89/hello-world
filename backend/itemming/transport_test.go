package itemming

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/itemming/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetItems(input model.ListInput) ([]model.Item, error) {
	args := m.Called(input)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *MockService) GetItem(uuid string) (model.Item, error) {
	args := m.Called(uuid)
	return args.Get(0).(model.Item), args.Error(1)
}

func (m *MockService) UpdateItem(item model.ItemInput) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockService) CreateItem(item model.ItemInput) (string, error) {
	args := m.Called(item)
	return args.String(0), args.Error(1)
}

func (m *MockService) CreateItems(items []model.ItemInput) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *MockService) DeleteItem(uuid string) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func (m *MockService) Health() utils.ServiceHealthCheck {
	args := m.Called()
	return args.Get(0).(utils.ServiceHealthCheck)
}

func TestRegisterRoutes(t *testing.T) {
	router := httprouter.New()
	mockService := new(MockService)
	RegisterRoutes(router, mockService)

	t.Run("GET /items/:uuid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items/test-uuid", nil)
		req = req.WithContext(context.Background())
		w := httptest.NewRecorder()

		mockService.On("GetItem", "mock-uuid").Return(GetItemResponse{item: model.Item{ID: "test-uuid"}}, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response GetItemResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "test-uuid", response.item.ID)
	})

	t.Run("GET /items/:uuid with error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items/nonexistent", nil)
		req = req.WithContext(context.Background())
		w := httptest.NewRecorder()

		mockService.On("GetItem", "nonexistent").Return(GetItemResponse{}, errors.New("not found"))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
