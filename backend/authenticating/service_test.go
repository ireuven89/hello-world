package authenticating

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ireuven89/hello-world/backend/authenticating/mocks"
	"github.com/ireuven89/hello-world/backend/authenticating/model"
)

func TestAuthService_RegisterSuccess(t *testing.T) {
	logger := zap.NewNop()
	mockInMemStore := mocks.InMemMock{Mock: mock.Mock{}}
	service := NewAuthService(&mockInMemStore, logger)

	//success mock
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

func TestAuthService_Login(t *testing.T) {
	mockUserStore := new(mocks.InMemMock)
	logger := zap.NewNop() // No-op logger to avoid real logging in tests
	service := &AuthService{
		userStore: mockUserStore,
		logger:    logger,
	}

	// Prepare test data
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	testUser := model.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}

	// Success Case: Valid username and password
	mockUserStore.Mock.On("Find", "testuser").Return(testUser, nil)

	token, err := service.Login("testuser", password)
	assert.NoError(t, err, "expected no error for valid login")
	assert.NotEmpty(t, token, "expected a valid JWT token")

	// Failure Case: User not found
	mockUserStore.Mock.On("Find", "unknown").Return(model.User{}, errors.New("user not found"))

	token, err = service.Login("unknown", password)
}
