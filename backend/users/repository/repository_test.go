package users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockRedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

type MockDB struct {
	mock.Mock
	*sqlz.DB
}

func (m *MockDB) Select(query string, args ...interface{}) *sqlz.SelectStmt {
	argsList := append([]interface{}{query}, args...)
	called := m.Called(argsList...)
	return called.Get(0).(*sqlz.SelectStmt)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	argsList := append([]interface{}{query}, args...)
	called := m.Called(argsList...)
	return called.Get(0).(sql.Result), called.Error(1)
}

func TestRepository_Get(t *testing.T) {
	mockDb := MockDB{}
	mockRedis := MockRedisClient{}
	logger := zaptest.NewLogger(t)

	repo := New(mockDb.DB, &mockRedis, logger)

	mockRedis.On("Get", fmt.Sprintf("FindUser:%s", "uuid")).Return(errors.New("not found"))
	res, err := repo.FindUser("uuid")

	assert.Nil(t, err)
	assert.NotEmpty(t, res)

}
