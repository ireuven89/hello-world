package biddering

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ireuven89/hello-world/backend/biddering/model"
)

type mockService struct{}

func (m *mockService) GetBidder(uuid string) (model.Bidder, error) {
	return model.Bidder{
		Id:          int64(1),
		Uuid:        uuid,
		UserUuid:    "user-123",
		Name:        "Test Bidder",
		Item:        "Test Item",
		Price:       100.0,
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *mockService) ListBidders(input model.BiddersInput) ([]model.Bidder, error) {
	return []model.Bidder{
		{
			Id:          int64(1),
			Uuid:        "uuid-1",
			UserUuid:    "user-123",
			Name:        "Bidder 1",
			Item:        "Item 1",
			Price:       50.0,
			Description: "Description 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Id:          int64(2),
			Uuid:        "uuid-2",
			UserUuid:    "user-456",
			Name:        "Bidder 2",
			Item:        "Item 2",
			Price:       150.0,
			Description: "Description 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

func (m *mockService) CreateBidder(req model.BidderInput) (string, error) {
	return "", nil
}

func (m *mockService) UpdateBidder(req model.BidderInput) (string, error) {
	return "", nil
}

func (m *mockService) Delete(id string) error {
	return nil
}

// Test RegisterRoutes
func TestRegisterRoutes(t *testing.T) {
	router := httprouter.New()
	service := &mockService{}

	RegisterRoutes(router, service)

	t.Run("GetBidder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bidders/uuid-123", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := map[string]interface{}{
			"id":          int64(1),
			"uuid":        "",
			"userUuid":    "user-123",
			"name":        "Test Bidder",
			"itemming":    "Test Item",
			"price":       100.0,
			"description": "Test Description",
			"created_at":  mock.Anything,
			"updated_at":  mock.Anything,
		}

		var actualResponse map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("ListBidders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bidders", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("CreateBidder", func(t *testing.T) {
		newBidder := map[string]interface{}{
			"name":        "New Bidder",
			"itemming":    "New Item",
			"price":       200.0,
			"description": "New Description",
		}
		body, _ := json.Marshal(newBidder)

		req := httptest.NewRequest(http.MethodPost, "/bidders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
