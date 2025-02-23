package repository

import (
	"testing"

	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRepo_Find(t *testing.T) {
	logger := zap.NewNop()
	mockDb, mock, err := sqlmock.New()
	sqlzMock := sqlz.New(mockDb, "mysql")
	user := "user"
	password := "password"
	repo := New(logger, sqlzMock)

	mockResult := sqlmock.NewRows(
		[]string{"user", "password"}).AddRow(user, password)
	mock.ExpectQuery("SELECT user, password FROM userring WHERE user = ?").WithArgs(user).WillReturnRows(mockResult)

	result, err := repo.Find(user)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotEmpty(t, result)

}

func TestRepo_Save(t *testing.T) {
	logger := zap.NewNop()
	mockDb, mock, err := sqlmock.New()
	sqlzMock := sqlz.New(mockDb, "mysql")
	user := "user"
	password := "password"
	repo := New(logger, sqlzMock)

	mock.ExpectExec("INSERT INTO userring ").
		WithArgs(sqlmock.AnyArg(), password, user). // Matching any argument for ID
		WillReturnResult(
<<<<<<< Updated upstream
			sqlmock.NewResult(1, 1), // Return the mock UUID as the inserted id
=======
			sqlmock.NewResult(1, 1), // Return the mocks UUID as the inserted id
>>>>>>> Stashed changes
		)

	err = repo.Save(user, password)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}
