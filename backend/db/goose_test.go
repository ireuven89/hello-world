package db

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"testing"
)

type MockGooseService struct {
	mock mock.Mock
}

func (mg *MockGooseService) lockDB() (*sql.Tx, error) {
	args := mg.mock.Called()

	res, ok := args.Get(0).(*sql.Tx)

	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}

func (mg *MockGooseService) unlockDB(transaction *sql.Tx) error {
	args := mg.mock.Called(transaction)

	return args.Error(0)
}

func (mg *MockGooseService) migrateDB() error {
	args := mg.mock.Called()

	return args.Error(0)
}

func TestLockingDB(t *testing.T) {
	mgs := MockGooseService{mock: mock.Mock{}}

	mgs.mock.On("lockDB").Return(nil, errors.New("failed locking db"))

	tx, err := mgs.lockDB()

	assert.NotNil(t, err)
	assert.Nil(t, tx)

}

func TestUnlockingDB(t *testing.T) {
	mgs := MockGooseService{mock: mock.Mock{}}
	tx := sql.Tx{}

	mgs.mock.On("unlockDB", &tx).Return(nil)

	err := mgs.unlockDB(&tx)

	assert.Nil(t, err)
}

func TestMigrateDB(t *testing.T) {
	mgs := MockGooseService{mock: mock.Mock{}}

	mgs.mock.On("migrateDB").Return(nil)

	err := mgs.migrateDB()

	assert.Nil(t, err)
}
