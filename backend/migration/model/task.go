package model

type MigrationTask struct {
	Name     string
	Execute  func() error
	Rollback func() error
}
