package userring

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/userring/mocks"
	"github.com/ireuven89/hello-world/backend/userring/model"
)

type CreateUserTest struct {
	name               string
	mockRepositoryCall *mock.Call
	input              model.UserUpsertInput
	wantErr            bool
	expected           string
}

func TestService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	logger := zap.NewNop()
	service := New(logger, mockRepo)
	testInputs := map[string]model.UserUpsertInput{
		"success": {
			Name: "name",
			Uuid: "mock-uuid",
		},
		"failure": {
			Name: "name",
			Uuid: "",
		},
	}

	tests := []CreateUserTest{
		{
			name:               "success",
			mockRepositoryCall: mockRepo.On("Upsert", testInputs["success"]).Return("mock-uuid", nil),
			input:              testInputs["success"],
			wantErr:            false,
			expected:           "mock-uuid",
		},
		{
			name:               "failure",
			mockRepositoryCall: mockRepo.On("Upsert", testInputs["failure"]).Return("", errors.New("database error")),
			input:              testInputs["failure"],
			wantErr:            true,
			expected:           "",
		},
	}

	for _, test := range tests {
		res, err := service.CreateUser(test.input)
		assert.Equal(t, err != nil, test.wantErr, test.name)
		assert.Equal(t, res, test.expected, test.name)
	}

	mockRepo.AssertExpectations(t) // Assert that Upsert was called as expected
}

type DeleteUserTest struct {
	name               string
	mockRepositoryCall *mock.Call
	input              string
	wantErr            bool
}

func TestService_DeleteUser(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	logger := zap.NewNop()
	service := New(logger, mockRepo)
	tests := []DeleteUserTest{
		{
			name:               "fail on invalid user",
			mockRepositoryCall: mockRepo.On("Delete", "").Return(errors.New("database error")),
			input:              "",
			wantErr:            true,
		},
		{
			name:               "success",
			mockRepositoryCall: mockRepo.On("Delete", "mock-uuid").Return(nil),
			input:              "mock-uuid",
			wantErr:            false,
		},
	}

	for _, test := range tests {
		err := service.DeleteUser(test.input)
		assert.Equal(t, err != nil, test.wantErr)
	}

	mockRepo.AssertExpectations(t) // Assert that Upsert was called as expected

}

type GetUserTest struct {
	name               string
	mockRepositoryCall *mock.Call
	input              string
	wantErr            bool
	expected           model.User
}

func TestService_GetUser(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	logger := zap.NewNop()
	service := New(logger, mockRepo)

	tests := []GetUserTest{
		{
			name:               "fail on invalid input",
			input:              "",
			mockRepositoryCall: mockRepo.On("FindUser", "").Return(model.User{}, errors.New("model not found")),
			wantErr:            true,
			expected:           model.User{},
		},
		{
			name:               "success",
			input:              "uuid",
			mockRepositoryCall: mockRepo.On("FindUser", "uuid").Return(model.User{Uuid: "uuid"}, nil),
			wantErr:            false,
			expected:           model.User{Uuid: "uuid"},
		},
	}

	for _, test := range tests {
		res, err := service.GetUser(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)

	}
}

func TestGoRoutine(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)       // Increment counter
		go task(i, &wg) // Run goroutine
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All tasks completed")
}

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement WaitGroup counter when done
	fmt.Printf("Task %d is running\n", id)
}
