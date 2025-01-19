package users

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ireuven89/hello-world/backend/users/model"
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
	sqlz *sqlz.DB
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

func TestRepository_FindUserWithoutCaching(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	mockSqlz := MockDB{sqlz: sqlz.New(mockdb, "mysql")}
	mockRedis := new(MockRedisClient)
	logger := zaptest.NewLogger(t)

	repo := New(mockSqlz.sqlz, mockRedis, logger)

	cachedQuery := fmt.Sprintf("FindUser:%s", "uuid")
	expectedQuery := `SELECT id, uuid, name, region FROM users WHERE uuid = ?`
	rows := sqlmock.NewRows([]string{"id", "uuid", "name", "region"}).
		AddRow(1, "1234", "John", "US")
	expectedResult := model.User{
		ID:     1,
		Uuid:   "1234",
		Name:   "John",
		Region: "US",
	}

	mock.ExpectQuery(expectedQuery).WithArgs("uuid").WillReturnRows(rows)

	// Step 4: Run the function being tested
	var result model.User
	mockRedis.On("Get", fmt.Sprintf("FindUser:%s", "uuid")).Return(nil, errors.New("not found"))
	mockRedis.On("Set", cachedQuery, expectedResult, redisQueryTTl).Return(nil)

	result, err = repo.FindUser("uuid")

	mockRedis.AssertCalled(t, "Get", cachedQuery)
	mockRedis.AssertCalled(t, "Set", cachedQuery, expectedResult, redisQueryTTl)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, expectedResult, result)

}

func TestRepository_GetWithCaching(t *testing.T) {
	mockdb, mockSql, err := sqlmock.New()
	mockSqlz := MockDB{sqlz: sqlz.New(mockdb, "mysql")}
	mockRedis := new(MockRedisClient)
	logger := zaptest.NewLogger(t)

	repo := New(mockSqlz.sqlz, mockRedis, logger)

	cachedQuery := fmt.Sprintf("FindUser:%s", "uuid")
	expectedQuery := `SELECT id, uuid, name, region FROM users WHERE uuid = ?`
	rows := sqlmock.NewRows([]string{"id", "uuid", "name", "region"}).
		AddRow(1, "1234", "John", "US")
	cachedUser := model.User{
		ID:     1,
		Uuid:   "1234",
		Name:   "John",
		Region: "US",
	}

	mockSql.ExpectQuery(expectedQuery).WithArgs("uuid").WillReturnRows(rows)

	// Step 4: Run the function being tested
	var result model.User
	mockRedis.On("Get", fmt.Sprintf("FindUser:%s", "uuid")).Return(cachedUser, nil)

	result, err = repo.FindUser("uuid")
	assert.NoError(t, err, "Error should be nil on cache hit")
	assert.Equal(t, cachedUser, result, "Returned model should match cached model")
	mockRedis.AssertCalled(t, "Get", cachedQuery) // Ensure cache is checked
	mockSql.ExpectationsWereMet()
}

func TestUserRepository_ListUsersWithCaching(t *testing.T) {
	mockdb, mockSql, err := sqlmock.New()
	mockSqlz := MockDB{sqlz: sqlz.New(mockdb, "mysql")}
	mockRedis := new(MockRedisClient)
	logger := zaptest.NewLogger(t)
	input := model.UserFetchInput{Name: "name"}

	repo := New(mockSqlz.sqlz, mockRedis, logger)

	cachedQuery := fmt.Sprintf("ListUsers:%s%s%s%v%v", input.Region, input.Name, input.Uuid, input.Page, input.Size)
	expectedQuery := `SELECT id, uuid, name, region FROM users WHERE name = ?`
	rows := sqlmock.NewRows([]string{"id", "uuid", "name", "region"}).
		AddRow(1, "1234", "name", "US")
	cachedUser := []model.User{{
		ID:     1,
		Uuid:   "1234",
		Name:   "name",
		Region: "US",
	},
	}

	mockSql.ExpectQuery(expectedQuery).WithArgs("name").WillReturnRows(rows)

	// Step 4: Run the function being tested
	var result []model.User
	mockRedis.On("Get", fmt.Sprintf("ListUsers:%s%s%s%v%v", input.Region, input.Name, input.Uuid, input.Page, input.Size)).Return(cachedUser, nil)

	result, err = repo.ListUsers(input)
	assert.NoError(t, err, "Error should be nil on cache hit")
	assert.Equal(t, cachedUser, result, "Returned model should match cached model")
	mockRedis.AssertCalled(t, "Get", cachedQuery) // Ensure cache is checked
	mockSql.ExpectationsWereMet()
}

func TestUpsert_Create(t *testing.T) {
	// Setup mocks database and logger
	logger := zap.NewNop() // No-op logger for testing purposes
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize mocks DB: %v", err)
		defer mockDb.Close()

		// Create a UserRepository instance
		repo := &UserRepository{
			db:     sqlz.New(mockDb, "mysql"),
			logger: logger,
		}

		// Input data for the Upsert method
		input := model.UserUpsertInput{
			Name:   "John Doe",
			Region: "North",
		}

		mockUuid := uuid.New().String()

		// Setup the expectation for the insert query
		mock.ExpectQuery("INSERT INTO users").WithArgs(sqlmock.AnyArg(), input.Name, input.Region).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(mockUuid),
		)

		// Run the Upsert method
		id, err := repo.Upsert(input)

		// Assertions
		assert.Nil(t, err)
		assert.Equal(t, mockUuid, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestUserRepository_Upsert_Update(t *testing.T) {
	// Initialize the logger and mocks database
	logger := zap.NewNop()
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize mocks DB: %v", err)
	// Initialize the logger and mock database
	defer mockDb.Close()

	// Create a UserRepository instance
	repo := &UserRepository{
		db:     sqlz.New(mockDb, "mysql"),
		logger: logger,
	}

	// Input data for the update scenario
	input := model.UserUpsertInput{
		Uuid:   "existing-uuid", // This triggers the update path
		Name:   "Updated Name",
		Region: "Updated Region",
	}

	// Mock the expected update query
	mock.ExpectQuery(`UPDATE users SET name = \?, region = \? WHERE uuid = \? RETURNING id`).
		WithArgs(input.Name, input.Region, input.Uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(input.Uuid))

	// Run the Upsert method
	id, err := repo.Upsert(input)

	// Assertions
	assert.Nil(t, err)                            // Ensure no error occurred
	assert.Equal(t, input.Uuid, id)               // Ensure the returned ID matches the UUID
}
