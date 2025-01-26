package authenticating

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeRegisterRequest(t *testing.T) {
	t.Run("Valid Request", func(t *testing.T) {
		// Prepare input
		input := RegisterRequest{
			UserName: "testuser",
			Password: "testpassword",
		}
		body, _ := json.Marshal(input)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		ctx := context.Background()

		// Call the function
		decodedRequest, err := decodeRegisterRequest(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, decodedRequest)

		result, ok := decodedRequest.(RegisterRequest)
		assert.True(t, ok)
		assert.Equal(t, input.UserName, result.UserName)
		assert.Equal(t, input.Password, result.Password)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		// Invalid JSON body
		body := []byte(`{"user": "testuser", "password": }`) // Malformed JSON

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		ctx := context.Background()

		// Call the function
		decodedRequest, err := decodeRegisterRequest(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, decodedRequest)
	})
}

func TestDecodeLoginRequest(t *testing.T) {
	t.Run("Valid Request", func(t *testing.T) {
		// Prepare valid input
		input := LoginRequestModel{
			UserName: "testuser",
			Password: "testpassword",
		}
		body, _ := json.Marshal(input)

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		ctx := context.Background()

		// Call the function
		decodedRequest, err := decodeLoginRequest(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, decodedRequest)

		result, ok := decodedRequest.(LoginRequestModel)
		assert.True(t, ok)
		assert.Equal(t, input.UserName, result.UserName)
		assert.Equal(t, input.Password, result.Password)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		// Prepare invalid JSON body
		body := []byte(`{"user": "testuser", "password": }`) // Malformed JSON

		// Create HTTP request
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		ctx := context.Background()

		// Call the function
		decodedRequest, err := decodeLoginRequest(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, decodedRequest)
	})
}

func TestEncodeLoginResponse(t *testing.T) {
	t.Run("Valid Response", func(t *testing.T) {
		// Prepare valid response
		response := LoginResponseModel{
			Token: "valid_token_123",
		}

		// Create HTTP ResponseRecorder
		rr := httptest.NewRecorder()
		ctx := context.Background()

		// Call the function
		err := encodeLoginResponse(ctx, rr, response)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		// Parse the response body
		var result map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "valid_token_123", result["token"])
	})

	t.Run("Invalid Response Type", func(t *testing.T) {
		// Prepare invalid response
		invalidResponse := "invalid_response"

		// Create HTTP ResponseRecorder
		rr := httptest.NewRecorder()
		ctx := context.Background()

		// Call the function
		err := encodeLoginResponse(ctx, rr, invalidResponse)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "encodeLoginResponse failed to decode response", err.Error())
	})
}

func TestDecodeVerifyRequest(t *testing.T) {
	t.Run("Valid Request", func(t *testing.T) {
		// Create a valid request body
		body := VerifyRequestModel{
			JwtToken: "valid_token_123",
		}
		bodyBytes, err := json.Marshal(body)
		assert.NoError(t, err)

		// Create an HTTP request
		req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Call the function
		ctx := context.Background()
		decodedReq, err := decodeVerifyRequest(ctx, req)

		// Assertions
		assert.NoError(t, err)

		verifyReq, ok := decodedReq.(VerifyRequestModel)
		assert.True(t, ok)
		assert.Equal(t, "valid_token_123", verifyReq.JwtToken)
	})

	t.Run("Invalid JSON Request", func(t *testing.T) {
		// Create an invalid JSON body
		body := `{"JwtToken":`
		req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		// Call the function
		ctx := context.Background()
		decodedReq, err := decodeVerifyRequest(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, decodedReq)
	})
}

func TestEncodeVerifyResponse(t *testing.T) {
	t.Run("Valid Response", func(t *testing.T) {
		// Create a valid VerifyResponseModel
		response := VerifyResponseModel{
			User: "test_user",
		}

		// Create a response writer
		writer := httptest.NewRecorder()

		// Call the function
		ctx := context.Background()
		err := encodeVerifyResponse(ctx, writer, response)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
		assert.Equal(t, "application/json", writer.Header().Get("Content-Type"))

		// Decode the response body
		var responseBody map[string]interface{}
		err = json.NewDecoder(writer.Body).Decode(&responseBody)
		assert.NoError(t, err)

		// Validate the response
		assert.Equal(t, "test_user", responseBody["user"])
	})

	t.Run("Invalid Response Type", func(t *testing.T) {
		// Create an invalid response type
		invalidResponse := "invalid"

		// Create a response writer
		writer := httptest.NewRecorder()

		// Call the function
		ctx := context.Background()
		err := encodeVerifyResponse(ctx, writer, invalidResponse)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "encodeVerifyResponse failed to decode response", err.Error())
	})
}
