package aws

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
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

func TestClient_DeleteObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}

	client.mock.On("DeleteObject", "mock-key", "mock-bucket").Return(nil)

	err := client.DeleteObject("mock-key", "mock-bucket")

	assert.Nil(t, err)
}

func TestClient_DeleteObjectFail(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}

	client.mock.On("DeleteObject", "mock-key", "").Return(errors.New("error deleting object - bucket not exists"))

	err := client.DeleteObject("mock-key", "")

	assert.NotNil(t, err)
}

func TestClient_GetObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	expectedResult := []byte{'t', 'h', 'e'}

	client.mock.On("GetObject", "mock-key", "mock-bucket").Return([]byte{'t', 'h', 'e'}, nil)

	res, err := client.GetObject("mock-key", "mock-bucket")

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, expectedResult, res)
}

func TestClient_GetObjectFail(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}

	client.mock.On("GetObject", "mock-key", "").Return([]byte{}, errors.New("bucket not found error"))

	res, err := client.GetObject("mock-key", "")

	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestClient_PutObjectFail(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	file, _ := os.Open("failed_file.exe")
	client.mock.On("PutObject", "mock-key", "mock-bucket", file).Return(errors.New("file not exits"))

	err := client.PutObject("mock-key", "mock-bucket", file)
	assert.NotNil(t, err)
}

func TestClient_PutObject(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	file, _ := os.Open("s3client.go")
	client.mock.On("PutObject", "mock-key", "mock-bucket", file).Return(nil)

	err := client.PutObject("mock-key", "mock-bucket", file)
	assert.Nil(t, err)
}
