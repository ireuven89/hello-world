package aws

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ireuven89/hello-world/backend/environment"
	"go.uber.org/zap"
)

type AWSClient interface {
	PutObject(key, bucket string, file os.File) error
	GetObject(key, bucket string) error
	DeleteObject(key, bucket string) error
}

type Client struct {
	client *s3.S3
	logger *zap.Logger
}

func New(logger *zap.Logger) (*Client, error) {
	config := aws.Config{Region: &environment.Variables.AwsRegion}
	var err error

	defer func() {
		if recover() != nil {
			err = errors.New("failed to initiate session")
		}
	}()
	sess := session.Must(session.NewSession(&config))
	client := s3.New(sess)

	return &Client{
		client: client,
		logger: logger,
	}, err
}

func (c *Client) PutObject(key, bucket string, file *os.File) error {
	input := s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(key), Body: file}
	out, err := c.client.PutObject(&input)

	if err != nil {
		return err
	}

	c.logger.Debug("uploaded successfully %v", zap.String("output", out.String()))

	return nil
}

func (c *Client) GetObject(key, bucket string) (interface{}, error) {
	input := s3.GetObjectInput{Bucket: &bucket, Key: &key}
	out, err := c.client.GetObject(&input)

	if err != nil {
		return nil, err
	}

	c.logger.Debug("downloaded file successfully %v", zap.String("output", out.String()))

	return out.Body, nil
}

func (c *Client) DeleteObject(key, bucket string) error {
	input := s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	res, err := c.client.DeleteObject(&input)

	if err != nil {
		return err
	}

	c.logger.Debug("deleted object %v ", zap.String("output", res.String()))

	return nil
}
