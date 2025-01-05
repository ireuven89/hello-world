package sns

import (
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"context"
	"go.uber.org/zap"
)

type SnsClient struct {
	client sns.Client
	logger *zap.Logger
}

// Publish - this method publishes a message to Sns
func (sc *SnsClient) Publish(ctx context.Context, message string) (*sns.PublishOutput, error) {

	output, err := sc.client.Publish(ctx, &sns.PublishInput{Message: &message})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (sc *SnsClient) Subscribe() error {

	return nil
}
