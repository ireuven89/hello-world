package users

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/users/model"
)

type MockService struct {
	logger *zap.Logger
	repo   UserRepository
	mock   mock.Mock
}

type MockUserRepository struct {
	mock mock.Mock
}

func (m *MockUserRepository) ListUsers(input model.UserFetchInput) ([]model.User, error) {
	args := m.mock.Called(input)

	return args.Get(0).([]model.User), args.Error(1)
}
func (m *MockUserRepository) FindUser(uuid string) (model.User, error) {
	args := m.mock.Called(uuid)

	return args.Get(0).(model.User), args.Error(1)
}
func (m *MockUserRepository) Upsert(input model.UserUpsertInput) (string, error) {
	args := m.mock.Called(input)

	return args.Get(0).(string), args.Error(1)
}
func (m *MockUserRepository) Delete(uuid string) error {
	args := m.mock.Called(uuid)

	return args.Error(0)
}

func (ms *MockService) CreateUser(input model.UserUpsertInput) (string, error) {
	args := ms.mock.Called(input)

	return args.Get(0).(string), args.Error(1)
}

func (ms *MockService) DeleteUser(uuid string) error {
	args := ms.mock.Called(uuid)

	return args.Error(0)
}

func (ms *MockService) GetUser(uuid string) (model.User, error) {
	args := ms.mock.Called(uuid)

	return args.Get(0).(model.User), args.Error(1)
}

type CreateTest struct {
	name     string
	mockCall *mock.Call
	input    model.UserUpsertInput
	wantErr  bool
	output   string
}

func TestService_CreateUser(t *testing.T) {
	userRepo := MockUserRepository{mock.Mock{}}
	mockService := MockService{mock: mock.Mock{}, repo: &userRepo}
	tests := []CreateTest{
		{
			name: "fail on invalid input",
			input: model.UserUpsertInput{
				Uuid: "",
				Name: "",
			},
			mockCall: mockService.mock.On("CreateUser", model.UserUpsertInput{
				Uuid: "",
				Name: "",
			}).Return("", errors.New("failed to create model missing input")),
			wantErr: true,
			output:  "",
		},
		{
			name: "success",
			input: model.UserUpsertInput{
				Uuid: "mock-uuid",
				Name: "mock-name",
			},
			mockCall: mockService.mock.On("CreateUser", model.UserUpsertInput{
				Uuid: "mock-uuid",
				Name: "mock-name",
			}).Return("mock-uuid", nil),
			wantErr: false,
			output:  "mock-uuid",
		},
	}

	for _, test := range tests {
		res, err := mockService.CreateUser(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, test.output, res)
	}

}

type DeleteTest struct {
	name     string
	mockCall *mock.Call
	input    string
	wantErr  bool
}

func TestService_DeleteUser(t *testing.T) {
	userRepo := MockUserRepository{mock.Mock{}}
	mockService := MockService{mock: mock.Mock{}, repo: &userRepo}
	tests := []DeleteTest{
		{
			name:     "fail on invalid input",
			input:    "",
			mockCall: mockService.mock.On("DeleteUser", "").Return(errors.New("failed to delete model not found")),
			wantErr:  true,
		},
		{
			name:     "success",
			input:    "uuid",
			mockCall: mockService.mock.On("DeleteUser", "uuid").Return(nil),
			wantErr:  false,
		},
	}

	for _, test := range tests {
		err := mockService.DeleteUser(test.input)
		assert.Equal(t, err != nil, test.wantErr)
	}
}

type GetUserTest struct {
	name               string
	mockServiceCall    *mock.Call
	mockRepositoryCall *mock.Call
	input              string
	wantErr            bool
	expected           model.User
}

func TestService_GetUser(t *testing.T) {
	mockRepo := MockUserRepository{mock: mock.Mock{}}
	mockService := MockService{mock: mock.Mock{}, repo: &mockRepo}
	tests := []GetUserTest{
		{
			name:               "fail on invalid input",
			input:              "",
			mockServiceCall:    mockService.mock.On("GetUser", "").Return(model.User{}, errors.New("failed to get model - not found")),
			mockRepositoryCall: mockRepo.mock.On("FindUser", "").Return(nil, errors.New("model not found")),
			wantErr:            true,
			expected:           model.User{},
		},
		{
			name:               "success",
			input:              "uuid",
			mockServiceCall:    mockService.mock.On("GetUser", "uuid").Return(model.User{Uuid: "uuid"}, nil),
			mockRepositoryCall: mockRepo.mock.On("FindUser", "uuid").Return(model.User{Uuid: "uuid"}, nil),
			wantErr:            false,
			expected:           model.User{Uuid: "uuid"},
		},
	}

	for _, test := range tests {
		res, err := mockService.GetUser(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)

	}
}
