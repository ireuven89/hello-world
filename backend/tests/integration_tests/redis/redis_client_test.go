package redis

import (
	"os"
	"testing"
	"time"

	"github.com/ireuven89/hello-world/backend/redis"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	key   = "test-key"
	value = "test-value"
	ttl   = time.Duration(3)
)

var redisClient *redis.Service

func init() {
	if err := os.Setenv("REDIS_HOST", "http://localhost/6379"); err != nil {
		panic(err)
	}
	logger := zap.New(zapcore.NewNopCore())

	c, err := redis.New(logger)
	if err != nil {
		panic("failed initializing service")
	}

	redisClient = c
}

func TestSet(t *testing.T) {
	err := redisClient.Set(key, value)

	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	val, err := redisClient.Get(key)

	assert.Nil(t, err)
	assert.NotEmpty(t, val)
	assert.Equal(t, val, value)
}
