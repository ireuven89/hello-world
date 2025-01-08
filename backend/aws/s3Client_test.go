package aws

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock mock.Mock
}

func (mc *MockClient) DeleteObject(key, bucket string) error {
	args := mc.mock.Called(key, bucket)

	return args.Error(0)
}

func (mc *MockClient) PutObject(key, bucket string, file *os.File) error {
	args := mc.mock.Called(key, bucket, file)

	return args.Error(0)
}

func (mc *MockClient) GetObject(key, bucket string) (interface{}, error) {
	args := mc.mock.Called(key, bucket)

	return args.Get(0).(interface{}), args.Error(1)
}

type deleteObjectTest struct {
	name    string
	wantErr bool
	input   struct {
		key    string
		bucket string
	}
	mockCall *mock.Call
}

func TestClient_DeleteObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	tests := []deleteObjectTest{
		{
			name:    "success",
			wantErr: false,
			input: struct {
				key    string
				bucket string
			}{key: "mock-key", bucket: "mock-bucket"},
			mockCall: client.mock.On("DeleteObject", "mock-key", "mock-bucket").Return(nil),
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mock-bucket"},
			mockCall: client.mock.On("DeleteObject", "", "mock-bucket").Return(errors.New("invalid key for s3")),
		},
		{
			name:    "fail invalid bucket",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mock-bucket"},
			mockCall: client.mock.On("DeleteObject", "mock-key", "").Return(errors.New("invalid bucket name")),
		},
	}

	for _, test := range tests {
		err := client.DeleteObject(test.input.key, test.input.bucket)
		assert.Equal(t, err != nil, test.wantErr)
	}
}

type getObjectTest struct {
	name    string
	wantErr bool
	input   struct {
		key    string
		bucket string
	}
	mockCall *mock.Call
	expected interface{}
}

func TestClient_GetObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	tests := []getObjectTest{
		{
			name:    "success",
			wantErr: false,
			input: struct {
				key    string
				bucket string
			}{key: "mock-key", bucket: "mock-bucket"},
			mockCall: client.mock.On("GetObject", "mock-key", "mock-bucket").Return([]byte{'r', 'e', 's'}, nil),
			expected: []byte{'r', 'e', 's'},
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mock-bucket"},
			mockCall: client.mock.On("GetObject", "", "mock-bucket").Return([]byte{}, errors.New("invalid key for s3")),
			expected: []byte{},
		},
	}

	for _, test := range tests {
		res, err := client.GetObject(test.input.key, test.input.bucket)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)
	}
}

type putObjectTest struct {
	name    string
	wantErr bool
	input   struct {
		key    string
		bucket string
		file   *os.File
	}
	mockCall *mock.Call
	expected interface{}
}

func TestClient_PutObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	tests := []putObjectTest{
		{
			name:    "success",
			wantErr: false,
			input: struct {
				key    string
				bucket string
				file   *os.File
			}{key: "mock-key", bucket: "mock-bucket"},
			mockCall: client.mock.On("PutObject", "mock-key", "mock-bucket", &os.File{}).Return(nil),
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
				file   *os.File
			}{key: "", bucket: "mock-bucket"},
			mockCall: client.mock.On("PutObject", "", "mock-bucket", &os.File{}).Return(errors.New("invalid key for s3")),
		},
	}

	for _, test := range tests {
		err := client.PutObject(test.input.key, test.input.bucket, &os.File{})
		assert.Equal(t, err != nil, test.wantErr)
	}
}
