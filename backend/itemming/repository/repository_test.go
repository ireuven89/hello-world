package repository

import (
	"errors"
	"testing"

	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ireuven89/hello-world/backend/itemming/mocks"
	"github.com/ireuven89/hello-world/backend/itemming/model"
)

func TestItemRepository_BulkInsert(t *testing.T) {
	mockDb, mockSql, _ := sqlmock.New()
	sqlzMock := sqlz.New(mockDb, "mysql")
	logger := zaptest.NewLogger(t)
	redisMock := new(mocks.MockRedis)
	repository := New(sqlzMock, logger, redisMock)

	input := []model.ItemInput{
		{
			Name:        "name",
			Description: "description-1",
			Price:       0,
		},
		{
			Name:        "new-name",
			Description: "description-2",
			Price:       1,
		},
	}

	mockSql.ExpectBegin()

	mockSql.ExpectExec("INSERT INTO items ").WithArgs(
		sqlmock.AnyArg(), "name", "description-1", 0, sqlmock.AnyArg(), "new-name", "description-2", 1).WillReturnResult(sqlmock.NewResult(1, 2))
	mockSql.ExpectCommit()
	err := repository.BulkInsert(input)

	assert.NoError(t, err)
}

func TestItemRepository_BulkInsertFail(t *testing.T) {
	mockDb, mockSql, _ := sqlmock.New()
	sqlzMock := sqlz.New(mockDb, "mysql")
	logger := zaptest.NewLogger(t)
	redisMock := new(mocks.MockRedis)
	repository := New(sqlzMock, logger, redisMock)

	input := []model.ItemInput{
		{
			Name:        "name",
			Description: "description-1",
			Price:       0,
		},
		{
			Name:        "new-name",
			Description: "description-2",
			Price:       1,
		},
	}

	mockSql.ExpectBegin()

	mockSql.ExpectExec("INSERT INTO items ").WithArgs(
		sqlmock.AnyArg(), "name", "description-1", 0, sqlmock.AnyArg(), "new-name", "description-2", 1).WillReturnResult(sqlmock.NewResult(1, 2))

	mockSql.ExpectCommit().WillReturnError(errors.New("failed committing changes"))
	mockSql.ExpectRollback()
	err := repository.BulkInsert(input)

	assert.Error(t, err)
	mockSql.ExpectationsWereMet()
}
