package aws

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"go.uber.org/zap/zaptest"

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
			}{key: "mocks-key", bucket: "mocks-bucket"},
			mockCall: client.mock.On("DeleteObject", "mocks-key", "mocks-bucket").Return(nil),
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mocks-bucket"},
			mockCall: client.mock.On("DeleteObject", "", "mocks-bucket").Return(errors.New("invalid key for s3")),
		},
		{
			name:    "fail invalid bucket",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mocks-bucket"},
			mockCall: client.mock.On("DeleteObject", "mocks-key", "").Return(errors.New("invalid bucket name")),
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
			}{key: "mocks-key", bucket: "mocks-bucket"},
			mockCall: client.mock.On("GetObject", "mocks-key", "mocks-bucket").Return([]byte{'r', 'e', 's'}, nil),
			expected: []byte{'r', 'e', 's'},
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
			}{key: "", bucket: "mocks-bucket"},
			mockCall: client.mock.On("GetObject", "", "mocks-bucket").Return([]byte{}, errors.New("invalid key for s3")),
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
			}{key: "mocks-key", bucket: "mocks-bucket"},
			mockCall: client.mock.On("PutObject", "mocks-key", "mocks-bucket", &os.File{}).Return(nil),
		},
		{
			name:    "fail invalid key",
			wantErr: true,
			input: struct {
				key    string
				bucket string
				file   *os.File
			}{key: "", bucket: "mocks-bucket"},
			mockCall: client.mock.On("PutObject", "", "mocks-bucket", &os.File{}).Return(errors.New("invalid key for s3")),
		},
	}

	for _, test := range tests {
		err := client.PutObject(test.input.key, test.input.bucket, &os.File{})
		assert.Equal(t, err != nil, test.wantErr)
	}
}

// MockS3Client implements s3iface.S3API for mocking
type MockS3Client struct {
	s3iface.S3API
	PutObjectFunc func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// PutObject mocks the PutObject method
func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunc(input)
}

// TestPutObject tests the PutObject method
func TestPutObject(t *testing.T) {
	logger := zaptest.NewLogger(t) // Create a logger for testing

	// Create a mocks S3 client
	mockClient := &MockS3Client{
		PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			// Verify input parameters for the mocks
			assert.Equal(t, aws.String("test-bucket"), input.Bucket)
			assert.Equal(t, aws.String("test-key"), input.Key)
			assert.NotNil(t, input.Body)
			return &s3.PutObjectOutput{}, nil
		},
	}

	// Create the ServiceAws instance with the mocked client
	service := &ServiceAws{
		client: mockClient,
		logger: logger,
	}

	// Create a temporary file to simulate the input file
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name()) // Clean up the file after the test
	_, _ = file.Write([]byte("test data"))
	_, _ = file.Seek(0, 0) // Reset the file pointer

	// Call the PutObject method
	err = service.PutObject("test-key", "test-bucket", file)

	// Validate the result
	assert.NoError(t, err)
}
