package authenticating

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/authenticating/mocks"
)

func TestAuthService_RegisterSuccess(t *testing.T) {
	logger := zap.NewNop()
	mockInMemStore := new(mocks.InMemMock)
	service := NewAuthService(mockInMemStore, logger)

	//success mocks
	user := "model"
	password := "password"
	mockInMemStore.Mock.On("Save", user, mock.Anything).Return(nil)

	err := service.Register(user, password)

	assert.Nil(t, err)
}

func TestAuthService_RegisterFail(t *testing.T) {
	logger := zap.NewNop()
	mockInMemStore := new(mocks.InMemMock)
	service := NewAuthService(mockInMemStore, logger)

	//success mocks
	user := "model"
	password := "password"
	mockInMemStore.Mock.On("Save", user, mock.Anything).Return(errors.New("invalid password"))

	err := service.Register(user, password)

	assert.Error(t, err)
}
