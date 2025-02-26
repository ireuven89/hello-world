package bider

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ireuven89/hello-world/backend/bider/model"
)

type MockRedis struct {
	mock mock.Mock
}

func (mdb *MockRedis) Set(key string, value interface{}, ttl time.Duration) error {
	args := mdb.mock.Called(key, value, ttl)

	return args.Error(0)
}

func (mdb *MockRedis) Get(key string) (interface{}, error) {
	args := mdb.mock.Called(key)

	return args.Get(0), args.Error(1)
}

func (mdb *MockRedis) Delete(uuid string) error {
	args := mdb.mock.Called(uuid)

	return args.Error(0)
}

func TestRepository_List(t *testing.T) {
	logger := zap.NewNop()
	redisMock := new(MockRedis)
	mockDB, mockSql, err := sqlmock.New()
	assert.NoError(t, err)
	mockSqlz := sqlz.New(mockDB, "mysql")
	createAt := time.Now()
	updateAt := time.Now()

	repo := New(mockSqlz, logger, redisMock)
	input := model.BiddersInput{
		Page: model.PageRequest{Offset: 0},
		Name: "name",
		Item: "itemming",
	}
	expectedResult := []model.Bidder{
		{
			Uuid:      "mock-uuid",
			Name:      input.Name,
			Item:      input.Item,
			CreatedAt: createAt,
			UpdatedAt: updateAt,
		},
	}

	redisQuery := fmt.Sprintf("%s%s%s%v%v", input.Uuid, input.Name, input.Item, input.Page.Offset, input.Page.GetLimit())
	rows := sqlmock.NewRows([]string{"uuid", "name", "itemming", "created_at", "updated_at"}).
		AddRow("mock-uuid", input.Name, input.Item, createAt, updateAt)
	mockSql.ExpectQuery("SELECT uuid, name, itemming, created_at, updated_at FROM bidders").
		WithArgs("name", "itemming").
		WillReturnRows(rows)

	//redis cache miss and set valid
	redisMock.mock.On("Get", redisQuery).Return(nil, errors.New("not found"))
	redisMock.mock.On("Set", redisQuery, expectedResult, redisQueryTtl).Return(nil)

	res, err := repo.List(input)

	redisMock.mock.AssertCalled(t, "Get", redisQuery)
	redisMock.mock.AssertCalled(t, "Set", redisQuery, expectedResult, redisQueryTtl)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, res)
	assert.NoError(t, mockSql.ExpectationsWereMet())

}

func TestRepository_Single(t *testing.T) {
	logger := zap.NewNop()
	redisMock := MockRedis{mock: mock.Mock{}}
	mockDB, mockSql, err := sqlmock.New()
	assert.NoError(t, err)
	mockSqlz := sqlz.New(mockDB, "mysql")
	createAt := time.Now()
	updateAt := time.Now()
	repo := New(mockSqlz, logger, &redisMock)
	mockUuid := "mock-uuid"
	input := model.BiddersInput{
		Page: model.PageRequest{Offset: 0},
		Name: "name",
		Item: "itemming",
	}
	expectedResult := model.Bidder{
		Uuid:      "mock-uuid",
		Name:      input.Name,
		Item:      input.Item,
		CreatedAt: createAt,
		UpdatedAt: updateAt,
	}
	rows := sqlmock.NewRows([]string{"uuid", "name", "itemming", "created_at", "updated_at"}).
		AddRow(mockUuid, "name", "itemming", createAt, updateAt)
	mockSql.ExpectQuery("SELECT uuid, name, itemming, created_at, updated_at FROM bidders").
		WithArgs(mockUuid).
		WillReturnRows(rows)

	res, err := repo.Single(mockUuid)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, res)
	assert.NoError(t, mockSql.ExpectationsWereMet())

}
