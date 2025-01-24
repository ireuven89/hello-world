package biddering

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ireuven89/hello-world/backend/biddering/mocks"
	"github.com/ireuven89/hello-world/backend/biddering/model"
)

func TestRepository_List(t *testing.T) {
	logger := zap.NewNop()
	redisMock := new(mocks.MockRedis)
	mockDB, mockSql, err := sqlmock.New()
	assert.NoError(t, err)
	mockSqlz := sqlz.New(mockDB, "mysql")
	createAt := time.Now()
	updateAt := time.Now()

	repo := NewRepository(mockSqlz, logger, redisMock)
	input := model.BiddersInput{
		Page: model.PageRequest{Offset: 0},
		Name: "name",
		Item: "item",
	}
	expectedResult := []model.Bidder{
		{
			Uuid:      "mocks-uuid",
			Name:      input.Name,
			Item:      input.Item,
			CreatedAt: createAt,
			UpdatedAt: updateAt,
		},
	}

	redisQuery := fmt.Sprintf("%s%s%s%v%v", input.Uuid, input.Name, input.Item, input.Page.Offset, input.Page.GetLimit())
	rows := sqlmock.NewRows([]string{"uuid", "name", "item", "created_at", "updated_at"}).
		AddRow("mocks-uuid", input.Name, input.Item, createAt, updateAt)
	mockSql.ExpectQuery("SELECT uuid, name, item, created_at, updated_at FROM bidders").
		WithArgs("name", "item").
		WillReturnRows(rows)

	//redis cache miss and set valid
	redisMock.Mock.On("Get", redisQuery).Return(nil, errors.New("not found"))
	redisMock.Mock.On("Set", redisQuery, expectedResult, redisQueryTtl).Return(nil)

	res, err := repo.List(input)

	redisMock.Mock.AssertCalled(t, "Get", redisQuery)
	redisMock.Mock.AssertCalled(t, "Set", redisQuery, expectedResult, redisQueryTtl)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, res)
	assert.NoError(t, mockSql.ExpectationsWereMet())

}

func TestRepository_Single(t *testing.T) {
	logger := zap.NewNop()
	redisMock := new(mocks.MockRedis)
	mockDB, mockSql, err := sqlmock.New()
	assert.NoError(t, err)
	mockSqlz := sqlz.New(mockDB, "mysql")
	createAt := time.Now()
	updateAt := time.Now()
	repo := NewRepository(mockSqlz, logger, redisMock)
	mockUuid := "mocks-uuid"
	input := model.BiddersInput{
		Page: model.PageRequest{Offset: 0},
		Name: "name",
		Item: "item",
	}
	expectedResult := model.Bidder{
		Uuid:      "mocks-uuid",
		Name:      input.Name,
		Item:      input.Item,
		CreatedAt: createAt,
		UpdatedAt: updateAt,
	}
	rows := sqlmock.NewRows([]string{"uuid", "name", "item", "created_at", "updated_at"}).
		AddRow(mockUuid, "name", "item", createAt, updateAt)
	mockSql.ExpectQuery("SELECT uuid, name, item, created_at, updated_at FROM bidders").
		WithArgs(mockUuid).
		WillReturnRows(rows)

	res, err := repo.FindOne(mockUuid)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, res)
	assert.NoError(t, mockSql.ExpectationsWereMet())

}
