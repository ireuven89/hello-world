package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ido50/sqlz"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestLockDB(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	logger := zaptest.NewLogger(t) // Use zaptest for testing zap logging
	ms := &MigrationService{
		db:     sqlz.New(db, "mysql"),
		logger: logger,
	}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError string
	}{
		{
			name: "successful lock",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("LOCK TABLE lock_table WRITE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: "",
		},
		{
			name: "begin transaction error",
			mockSetup: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedError: "failed to lock DB begin error",
		},
		{
			name: "lock query execution error",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("LOCK TABLE lock_table WRITE").WillReturnError(errors.New("exec error"))
			},
			expectedError: "failed to lock DB exec error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			// Call the function under test
			tx, err := ms.lockDB()

			// Check the result
			if test.expectedError == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tx == nil {
					t.Errorf("expected transaction, got nil")
				}
			} else {
				if err == nil || err.Error() != test.expectedError {
					t.Errorf("expected error: %v, got: %v", test.expectedError, err)
				}
				if tx != nil {
					t.Errorf("expected no transaction, got: %v", tx)
				}
			}

			// Assert expectations
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestIsTableLocked(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlzMock := sqlz.New(db, "mysql")
	assert.NoError(t, err)
	defer db.Close()

	ms := &MigrationService{db: sqlzMock}

	// 🟢 **Test Case 1: Table is locked**
	mock.ExpectQuery("SHOW OPEN TABLES WHERE In_use > 0 AND Table_name = ?").
		WithArgs("lock_table").
		WillReturnRows(sqlmock.NewRows([]string{"Database", "Table", "In_use", "Is_locked"}).
			AddRow("test_db", "lock_table", 1, 0))

	locked, err := ms.isTableLocked()
	assert.NoError(t, err)
	assert.True(t, locked, "Expected table to be locked")

	mock.ExpectQuery("SHOW OPEN TABLES WHERE In_use > 0 AND Table_name = ?").
		WithArgs("lock_table").
		WillReturnRows(sqlmock.NewRows([]string{"Database", "Table", "In_use", "Is_locked"}))

	locked, err = ms.isTableLocked()
	assert.NoError(t, err)
	assert.False(t, locked, "Expected table to be unlocked")

	mock.ExpectQuery("SHOW OPEN TABLES WHERE In_use > 0 AND Table_name = ?").
		WithArgs("lock_table").
		WillReturnError(errors.New("database error"))

	locked, err = ms.isTableLocked()
	assert.Error(t, err)
	assert.False(t, locked, "Expected function to return false on error")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlockDB(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	logger := zaptest.NewLogger(t) // Use zaptest for testing zap logging
	ms := &MigrationService{
		db:     sqlz.New(db, "mysql"),
		logger: logger,
	}

	tests := []struct {
		name          string
		transaction   *sql.Tx
		mockSetup     func()
		expectedError string
	}{
		{
			name:        "successful unlock",
			transaction: mockTransaction(db, mock),
			mockSetup: func() {
				mock.ExpectExec("UNLOCK TABLES").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(nil)
			},
			expectedError: "",
		},
		{
			name:          "nil transaction",
			transaction:   nil,
			mockSetup:     func() {}, // No setup needed for nil transaction case
			expectedError: "transaction not exists can't unlock",
		},
		{
			name:        "unlock query execution error",
			transaction: mockTransaction(db, mock),
			mockSetup: func() {
				mock.ExpectExec("UNLOCK TABLES").WillReturnError(errors.New("exec error"))
			},
			expectedError: "exec error",
		},
		{
			name:        "commit error",
			transaction: mockTransaction(db, mock),
			mockSetup: func() {
				mock.ExpectExec("UNLOCK TABLES").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedError: "failed to commit transaction commit error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			// Call the function under test
			err := ms.unlockDB(test.transaction)

			// Check the result
			if test.expectedError == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil || err.Error() != test.expectedError {
					t.Errorf("expected error: %v, got: %v", test.expectedError, err)
				}
			}

			// Assert expectations
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

// mockTransaction is a helper function to create a mocks *sql.Tx
func mockTransaction(db *sql.DB, mock sqlmock.Sqlmock) *sql.Tx {
	// Start a transaction from the mocked DB
	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		panic("failed to create mocked transaction: " + err.Error())
	}
	return tx
}
