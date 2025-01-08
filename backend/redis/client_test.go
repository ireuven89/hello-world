package redis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type test struct {
	name  string
	input struct {
		key   string
		value interface{}
	}
	mockCall *mock.Call
	wantErr  bool
	expected interface{}
}

type mockService struct {
	mock mock.Mock
}

func (ms *mockService) Set(key string, value interface{}) error {
	args := ms.mock.Called(key, value)

	return args.Error(0)
}

func (ms *mockService) Get(key string) (interface{}, error) {
	args := ms.mock.Called(key)

	return args.Get(0), args.Error(1)
}

func TestService_Set(t *testing.T) {
	ms := mockService{mock: mock.Mock{}}

	tests := []test{
		{
			name:     "empty string key",
			wantErr:  true,
			mockCall: ms.mock.On("Set", "", "").Return(errors.New("failed setting redis key")),
			input: struct {
				key   string
				value interface{}
			}{key: "", value: ""},
		},
		{
			name:     "valid string key",
			wantErr:  false,
			mockCall: ms.mock.On("Set", "key", "value").Return(nil),
			input: struct {
				key   string
				value interface{}
			}{key: "key", value: "value"},
		},
	}

	for _, test := range tests {
		err := ms.Set(test.input.key, test.input.value)
		assert.Equal(t, err != nil, test.wantErr)
	}
}

func TestService_Get(t *testing.T) {
	ms := mockService{mock: mock.Mock{}}

	tests := []test{
		{
			name:     "key not exists",
			wantErr:  true,
			mockCall: ms.mock.On("Get", "").Return("", errors.New("key not found")),
			input: struct {
				key   string
				value interface{}
			}{key: "", value: ""},
			expected: "",
		},
		{
			name:     "key exists",
			wantErr:  false,
			mockCall: ms.mock.On("Get", "key").Return("result", nil),
			input: struct {
				key   string
				value interface{}
			}{key: "key", value: "value"},
			expected: "result",
		},
	}

	for _, test := range tests {
		actual, err := ms.Get(test.input.key)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, test.expected, actual)
	}
}
