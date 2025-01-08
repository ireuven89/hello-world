package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/ireuven89/hello-world/backend/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const testBucket = "namespaces-test"

var awsClient *aws.Client

func init() {
	if err := os.Setenv("AWS_REGION", "us-east-1"); err != nil {
		panic(fmt.Sprintf("failed to set vars %v", err))
	}

	logger := zap.New(zapcore.NewNopCore())
	c, err := aws.New(logger)
	if err != nil {
		panic(fmt.Sprintf("failed initialies client test %v", err))
	}

	awsClient = c
}

func TestS3ClientPut(t *testing.T) {
	file, err := os.Open("./files/put_test_file.json")

	assert.Nilf(t, err, "faild to open test file")

	err = awsClient.PutObject("test-file", testBucket, file)

	assert.Nilf(t, err, fmt.Sprintf("fail on putting object to test bucket %v", err))
}

func TestS3ClientGet(t *testing.T) {
	file, err := os.Open("./files/put_test_file.json")

	assert.Nilf(t, err, "faild to open test file")

	err = awsClient.PutObject("test-file", testBucket, file)

	assert.Nilf(t, err, fmt.Sprintf("fail on putting object to test bucket %v", err))

	//get
	res, err := awsClient.GetObject("test-file", testBucket)
	assert.Nilf(t, err, "failed getting test file")
	assert.NotEmpty(t, res, "failed getting test file")
	assert.Nil(t, err)
}

func TestS3ClientDelete(t *testing.T) {
	file, err := os.Open("./files/put_test_file.json")

	assert.Nilf(t, err, "faild to open test file")

	err = awsClient.PutObject("test-file", testBucket, file)

	assert.Nilf(t, err, fmt.Sprintf("fail on putting object to test bucket %v", err))

	//delete
	err = awsClient.DeleteObject("test-file", testBucket)
	assert.Nilf(t, err, "failed deleting test file")
}
