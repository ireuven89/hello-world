package db

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/ido50/sqlz"
	mock_db "github.com/ireuven89/hello-world/backend/db/mock"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

func getMockDB() *sqlz.DB {

	return &sqlz.DB{
		DB: sqlx.NewDb(&sql.DB{}, "mock"),
	}
}

type DB interface {
	Begin() (sqlx.Tx, error)
}

func TestLockingDB(t *testing.T) {
	ctrl := gomock.Controller{T: t}
	mdb := mock_db.NewMockMigrationService(&ctrl)

	mdb.EXPECT()
}

func TestUnlockDB(t *testing.T) {
	logger := zap.Logger{}
	goose := New(getMockDB(), logger)

	tx, err := goose.lockDB()

	assert.Nil(t, err)

	err = goose.unlockDB(tx)

	assert.Nil(t, err)
}

func TestMigration(t *testing.T) {

}
