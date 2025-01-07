package redis

import (
	"context"
	"fmt"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Redis interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type Service struct {
	client *redis.Client
	logger *zap.Logger
}

var ctx = context.Background()
var ttl = time.Second * 3

func New(logger *zap.Logger) (*Service, error) {
	host := environment.Variables.RedisHost
	client := redis.NewClient(&redis.Options{
		Addr: host,
		DB:   0,
	})

	//ping check
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error(fmt.Sprintf("failed connecting to redis %v", err))
		return nil, err
	}

	return &Service{
		client: client,
		logger: logger,
	}, nil
}

func (s *Service) Set(key string, value interface{}) error {

	if err := s.client.Set(ctx, key, value, ttl).Err(); err != nil {
		s.logger.Error(fmt.Sprintf("failed inserting to redis %v", err))
		return err
	}

	return nil
}

func (s *Service) Get(key string) (interface{}, error) {

	result, err := s.client.Get(ctx, key).Result()

	if err != nil {
		s.logger.Error("failed to get value from redis")
		return nil, err
	}

	return result, nil
}
