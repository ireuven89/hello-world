package db

import (
	"time"

	"github.com/ido50/sqlz"
	"github.com/pressly/goose/v3"

	"github.com/ireuven89/hello-world/backend/db/model"

	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

const (
	lockMode = "WRITE"
)

type Service interface {
	lockDB() (*sql.Tx, error)
	unlockDB(transaction *sql.Tx) error
	migrateDB() error
	Run() error
}

type MigrationService struct {
	db            *sqlz.DB
	logger        *zap.Logger
	migrationsDir string
}

func New(db *sqlz.DB, logger *zap.Logger, migrationsDir string) Service {

	return &MigrationService{
		db:            db,
		logger:        logger,
		migrationsDir: migrationsDir,
	}
}

// Run - this function migrates DB
func (ms *MigrationService) Run() error {
	tx, err := ms.lockDB()

	if err != nil {
		ms.logger.Error("failed locking DB:", zap.Error(err))
		return err
	}

	if err = ms.migrateDB(); err != nil {
		ms.logger.Error("failed locking DB", zap.Error(err))
		return err
	}

	if err = ms.unlockDB(tx); err != nil {
		ms.logger.Error("failed unlocking DB", zap.Error(err))
		return err
	}

	return nil
}

// lockDB - this function locks the DB
func (ms *MigrationService) lockDB() (*sql.Tx, error) {

	tx, err := ms.db.Begin()

	if err != nil {
		ms.logger.Error("failed locking DB", zap.Error(err))
		return nil, fmt.Errorf("failed to lock DB %v", err)
	}

	lockQuery := fmt.Sprintf("LOCK TABLE %s %s", model.LockTable, lockMode)
	_, err = tx.Exec(lockQuery)

	if err != nil {
		ms.logger.Error("failed locking DB", zap.Error(err))
		return nil, fmt.Errorf("failed to lock DB %v", err)
	}

	return tx, nil
}

// unlockDB - this function unlocks DB given a transaction
func (ms *MigrationService) unlockDB(transaction *sql.Tx) error {

	if transaction == nil {
		ms.logger.Error("failed unlocking DB: transaction not exists")
		return fmt.Errorf("transaction not exists can't unlock")
	}

	unlockQuery := fmt.Sprintf("UNLOCK TABLES")
	_, err := transaction.Exec(unlockQuery)

	if err != nil {
		return err
	}

	if err = transaction.Commit(); err != nil {
		ms.logger.Error("failed unlocking DB", zap.Error(err))
		return fmt.Errorf("failed to commit transaction %v", err)
	}

	return nil
}

// migrateDB - this function migrates the DB
func (ms *MigrationService) migrateDB() error {
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	start := time.Now()
	ms.logger.Info("starting migration...")
	if err := goose.Up(ms.db.DB.DB, ms.migrationsDir); err != nil {
		ms.logger.Error("failed migration")
		return err
	}
	end := time.Now()
	ms.logger.Info("finished migration time to run: ", zap.Any("seconds", end.Second()-start.Second()))

	return nil
}
