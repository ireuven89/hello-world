package biddering

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ireuven89/hello-world/backend/biddering/mocks"
	"github.com/ireuven89/hello-world/backend/biddering/model"
)

func TestMakeEndpointListBidders_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockService)
	input := model.BiddersInput{Item: "item"}
	expectedBidders := []model.Bidder{
		{Uuid: "1", Name: "Bidder1"},
		{Uuid: "2", Name: "Bidder2"},
	}
	mockService.On("ListBidders", input).Return(expectedBidders, nil)

	endpointFunc := MakeEndpointListBidders(mockService)
	req := ListBiddersRequest{Input: input}

	// Act
	response, err := endpointFunc(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	res, ok := response.(ListBiddersResponseModel)
	assert.True(t, ok)
	assert.Equal(t, expectedBidders, res.bidders)
	mockService.AssertCalled(t, "ListBidders", input)
}

func TestMakeEndpointListBidders_Error(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockService)
	input := model.BiddersInput{Item: "Invalid"}
	mockService.On("ListBidders", input).Return([]model.Bidder{}, errors.New("service error"))

	endpointFunc := MakeEndpointListBidders(mockService)
	req := ListBiddersRequest{Input: input}

	// Act
	response, err := endpointFunc(context.Background(), req)

	// Assert
	assert.Nil(t, response)
	assert.EqualError(t, err, "MakeEndpointListBidders: service error")
	mockService.AssertCalled(t, "ListBidders", input)
}

func TestMakeEndpointListBidders_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockService)

	endpointFunc := MakeEndpointListBidders(mockService)

	// Act
	response, err := endpointFunc(context.Background(), "invalid-request")

	// Assert
	assert.Nil(t, response)
	assert.EqualError(t, err, "MakeEndpointListBidders failed cast request")
}

func TestMakeEndpointDeleteBidder(t *testing.T) {
	mockService := new(mocks.MockService)

	// Mock behavior
	validUUID := "12345"
	mockService.On("Delete", validUUID).Return(nil)
	invalidUUID := "invalid-id"
	mockService.On("Delete", invalidUUID).Return(errors.New("not found"))
	mockService.On("Delete", "").Return(errors.New("not found"))

	endpoint := MakeEndpointDeleteBidder(mockService)

	tests := []struct {
		name      string
		request   DeleteBidderRequestModel
		wantError bool
	}{
		{
			name:      "Valid request",
			request:   DeleteBidderRequestModel{Uuid: validUUID},
			wantError: false,
		},
		{
			name:      "Invalid request - bidder not found",
			request:   DeleteBidderRequestModel{Uuid: invalidUUID},
			wantError: true,
		},
		{
			name: "Invalid request type",
			request: DeleteBidderRequestModel{
				Uuid: "",
			}, // Not a DeleteBidderRequestModel
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Call the endpoint
			response, err := endpoint(ctx, tt.request)

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Nil(t, response)
			} else {
				assert.Nil(t, err)
				assert.Nil(t, response)
			}
		})
	}

	// Ensure expectations were met
	mockService.AssertExpectations(t)
}

func TestMakeEndpointGetBidder(t *testing.T) {
	mockService := new(mocks.MockService)

	// Mock data
	validUUID := "12345"
	invalidUUID := "invalid-id"

	mockBidder := model.Bidder{
		Uuid: "12345",
		Name: "Test Bidder",
	}

	// Mock behavior
	mockService.On("GetBidder", validUUID).Return(mockBidder, nil)
	mockService.On("GetBidder", invalidUUID).Return(model.Bidder{}, errors.New("bidder not found"))

	endpoint := MakeEndpointGetBidder(mockService)

	tests := []struct {
		name      string
		request   interface{}
		wantError bool
		expected  model.Bidder
	}{
		{
			name:      "Valid request",
			request:   GetBidderRequest{uuid: validUUID},
			wantError: false,
			expected:  mockBidder,
		},
		{
			name:      "Invalid request - bidder not found",
			request:   GetBidderRequest{uuid: invalidUUID},
			wantError: true,
			expected:  model.Bidder{},
		},
		{
			name:      "Invalid request type",
			request:   nil, // Invalid request type
			wantError: true,
			expected:  model.Bidder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Call the endpoint
			response, err := endpoint(ctx, tt.request)

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Nil(t, response)
			} else {
				assert.Nil(t, err)

				res, ok := response.(model.Bidder)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, res)
			}
		})
	}

	// Ensure expectations were met
	mockService.AssertExpectations(t)
}

func TestMakeEndpointCreateBidder(t *testing.T) {
	mockService := new(mocks.MockService)

	// Mock input and output
	validInput := model.BidderInput{
		Name:  "Test Bidder",
		Price: "100.0",
	}

	invalidInput := model.BidderInput{}

	mockUUID := "12345-uuid"
	mockError := errors.New("failed to create bidder")

	// Mock behavior
	mockService.On("CreateBidder", validInput).Return(mockUUID, nil)
	mockService.On("CreateBidder", invalidInput).Return("", mockError)

	endpoint := MakeEndpointCreateBidder(mockService)

	tests := []struct {
		name      string
		request   interface{}
		wantError bool
		expected  CreateBidderResponseModel
	}{
		{
			name: "Valid request",
			request: CreateBidderRequestModel{
				Input: validInput,
			},
			wantError: false,
			expected: CreateBidderResponseModel{
				ID: mockUUID,
			},
		},
		{
			name: "Invalid request - service error",
			request: CreateBidderRequestModel{
				Input: invalidInput,
			},
			wantError: true,
			expected:  CreateBidderResponseModel{},
		},
		{
			name:      "Invalid request type",
			request:   nil, // Invalid request type
			wantError: true,
			expected:  CreateBidderResponseModel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Call the endpoint
			response, err := endpoint(ctx, tt.request)

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Nil(t, response)
			} else {
				assert.Nil(t, err)

				res, ok := response.(CreateBidderResponseModel)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, res)
			}
		})
	}

	// Ensure all expectations were met
	mockService.AssertExpectations(t)
}

func TestMakeEndpointUpdateBidder(t *testing.T) {
	mockService := new(mocks.MockService)

	// Mock input and behavior
	validInput := model.BidderInput{
		Name:  "Updated Bidder",
		Price: "200.0",
	}

	invalidInput := model.BidderInput{}

	mockUUID := "12345-uuid"
	mockError := errors.New("failed to update bidder")

	// Mock service responses
	mockService.On("UpdateBidder", validInput).Return(mockUUID, nil)
	mockService.On("UpdateBidder", invalidInput).Return("", mockError)

	endpoint := MakeEndpointUpdateBidder(mockService)

	tests := []struct {
		name      string
		request   interface{}
		wantError bool
	}{
		{
			name: "Valid request",
			request: UpdateBidderRequestModel{
				Input: validInput,
			},
			wantError: false,
		},
		{
			name: "Invalid request - service error",
			request: UpdateBidderRequestModel{
				Input: invalidInput,
			},
			wantError: true,
		},
		{
			name:      "Invalid request type",
			request:   nil, // Invalid request type
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Call the endpoint
			response, err := endpoint(ctx, tt.request)

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Nil(t, response)
			} else {
				assert.Nil(t, err)
				assert.Nil(t, response)
			}
		})
	}

	// Ensure all expectations were met
	mockService.AssertExpectations(t)
}
