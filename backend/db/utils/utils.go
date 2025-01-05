package utils

import (
	"github.com/ido50/sqlz"
	"go.uber.org/zap"
)

type Utils struct {
	logger *zap.Logger
}

func New() *Utils {
	logger, _ := zap.NewDevelopment()

	return &Utils{
		logger: logger,
	}
}

func (u *Utils) DebugSelect(q *sqlz.SelectStmt, queryName string) {

	u.logger.Debug(queryName, zap.Any("statement: ", q))
}

func (u *Utils) DebugUpdate(q *sqlz.UpdateStmt, queryName string) {

	u.logger.Debug(queryName, zap.Any("statement: ", q))
}
func (u *Utils) DebugInsert(q *sqlz.InsertStmt, queryName string) {

	u.logger.Debug(queryName, zap.Any("statement: ", q))
}
func (u *Utils) DebugDelete(q *sqlz.DeleteStmt, queryName string) {

	u.logger.Debug(queryName, zap.Any("statement", q))
}
