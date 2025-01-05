package db

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

func (u *Utils) DebugSelect(q *sqlz.SelectStmt, queryName string) error {

}
