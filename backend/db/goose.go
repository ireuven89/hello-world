package db

import (
	"github.com/ireuven89/hello-world/backend/db/model"
	"github.com/pressly/goose/v3"
	"time"

	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

const (
	lockMode = "WRITE"
)

type MigrationService interface {
	lockDB() (*sql.Tx, error)
	unlockDB(transaction *sql.Tx) error
	migrateDB(transaction *sql.Tx) error
	validateDB() error
}

type Service struct {
	db            *sql.DB
	logger        *zap.Logger
	migrationsDir string
}

func New(db *sql.DB, logger *zap.Logger, migrationsDir string) *Service { // setup database

	return &Service{
		db:            db,
		logger:        logger,
		migrationsDir: migrationsDir,
	}
}

// Run - this function migrates DB
func (ms *Service) Run() error {
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
func (ms *Service) lockDB() (*sql.Tx, error) {

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

// unlockDB - this function unlocks DB give na transaction
func (ms *Service) unlockDB(transaction *sql.Tx) error {

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
func (ms *Service) migrateDB() error {
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	start := time.Now()
	ms.logger.Info("starting migration...")
	if err := goose.Up(ms.db, ms.migrationsDir); err != nil {
		ms.logger.Error("failed migration")
		return err
	}
	end := time.Now()
	ms.logger.Info("finished migration time to run: ", zap.Any("seconds", end.Second()-start.Second()))

	return nil
}
