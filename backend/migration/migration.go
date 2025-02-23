package migration

import (
	"context"
	"time"

	"github.com/sethvargo/go-retry"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ireuven89/hello-world/backend/utils"
)

func MustNewDB(config utils.DataBaseConnection) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	uri := options.Client().ApplyURI(config.Host)

	maxRet := retry.WithMaxRetries(3, retry.NewConstant(2*time.Second))
	client, err := retry.DoValue(context.Background(), maxRet, func(ctx context.Context) (*mongo.Client, error) {
		client, err := mongo.Connect(ctx, uri)
		if err != nil {
			return nil, retry.RetryableError(err)
		}

		return client, nil
	})

	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
