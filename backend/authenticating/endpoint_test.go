package authenticating

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ireuven89/hello-world/backend/authenticating/mocks"
)

// TestMakeEndpointRegister tests the MakeEndpointRegister function
func TestMakeEndpointRegister(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("Register", "testUser", "testPassword").Return(nil)

	endpoint := MakeEndpointRegister(mockService)
	ctx := context.Background()

	request := RegisterRequest{
		UserName: "testUser",
		Password: "testPassword",
	}

	_, err := endpoint(ctx, request)
	assert.NoError(t, err)

	mockService.AssertCalled(t, "Register", "testUser", "testPassword")
}

// TestMakeEndpointRegister_Failed tests the failure case for MakeEndpointRegister
func TestMakeEndpointRegister_Failed(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("Register", "testUser", "testPassword").Return(errors.New("registration error"))

	endpoint := MakeEndpointRegister(mockService)
	ctx := context.Background()

	request := RegisterRequest{
		UserName: "testUser",
		Password: "testPassword",
	}

	_, err := endpoint(ctx, request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registration error")

	mockService.AssertCalled(t, "Register", "testUser", "testPassword")
}

// TestMakeEndpointLogin tests the MakeEndpointLogin function
func TestMakeEndpointLogin(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("Login", "testUser", "testPassword").Return("testToken", nil)

	endpoint := MakeEndpointLogin(mockService)
	ctx := context.Background()

	request := LoginRequestModel{
		UserName: "testUser",
		Password: "testPassword",
	}

	response, err := endpoint(ctx, request)
	assert.NoError(t, err)

	res, ok := response.(LoginResponseModel)
	assert.True(t, ok)
	assert.Equal(t, "testToken", res.Token)

	mockService.AssertCalled(t, "Login", "testUser", "testPassword")
}

// TestMakeEndpointLogin_Failed tests the failure case for MakeEndpointLogin
func TestMakeEndpointLogin_Failed(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("Login", "testUser", "testPassword").Return("", errors.New("login error"))

	endpoint := MakeEndpointLogin(mockService)
	ctx := context.Background()

	request := LoginRequestModel{
		UserName: "testUser",
		Password: "testPassword",
	}

	_, err := endpoint(ctx, request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "login error")

	mockService.AssertCalled(t, "Login", "testUser", "testPassword")
}

// TestMakeEndpointVerify tests the MakeEndpointVerify function
func TestMakeEndpointVerify(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("VerifyToken", "testJwtToken").Return("testUser", nil)

	endpoint := MakeEndpointVerify(mockService)
	ctx := context.Background()

	request := VerifyRequestModel{
		JwtToken: "testJwtToken",
	}

	response, err := endpoint(ctx, request)
	assert.NoError(t, err)

	res, ok := response.(VerifyResponseModel)
	assert.True(t, ok)
	assert.Equal(t, "testUser", res.User)

	mockService.AssertCalled(t, "VerifyToken", "testJwtToken")
}

// TestMakeEndpointVerify_Failed tests the failure case for MakeEndpointVerify
func TestMakeEndpointVerify_Failed(t *testing.T) {
	mockService := new(mocks.MockService)
	mockService.On("VerifyToken", "testJwtToken").Return("", errors.New("verification error"))

	endpoint := MakeEndpointVerify(mockService)
	ctx := context.Background()

	request := VerifyRequestModel{
		JwtToken: "testJwtToken",
	}

	_, err := endpoint(ctx, request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verification error")

	mockService.AssertCalled(t, "VerifyToken", "testJwtToken")
}
