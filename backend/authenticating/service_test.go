package authenticating

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/authenticating/model"
)

type InMemMock struct {
	mock mock.Mock
}

func (ma *InMemMock) Save(username, password string) error {
	args := ma.mock.Called(username, password)

	return args.Error(0)
}
func (ma *InMemMock) Find(username string) (model.User, error) {
	args := ma.mock.Called(username)

	return args.Get(0).(model.User), args.Error(1)
}

func (ma *InMemMock) FindAll(page model.Page) ([]model.User, error) {
	args := ma.mock.Called(page)

	return args.Get(0).([]model.User), args.Error(1)
}

func (ma *InMemMock) Delete(id string) error {
	args := ma.mock.Called(id)

	return args.Error(0)
}

func TestAuthService_RegisterSuccess(t *testing.T) {
	logger := zap.NewNop()
	mockInMemStore := InMemMock{mock: mock.Mock{}}
	service := NewAuthService(&mockInMemStore, logger)

	//success mock
	user := "model"
	password := "password"
	mockInMemStore.mock.On("Save", user, mock.Anything).Return(nil)

	err := service.Register(user, password)

	assert.Nil(t, err)
}

func TestAuthService_RegisterFail(t *testing.T) {
	logger := zap.NewNop()
	mockInMemStore := InMemMock{mock: mock.Mock{}}
	service := NewAuthService(&mockInMemStore, logger)

	//success mock
	user := "model"
	password := "password"
	mockInMemStore.mock.On("Save", user, mock.Anything).Return(errors.New("invalid password"))

	err := service.Register(user, password)

	assert.Error(t, err)
}
