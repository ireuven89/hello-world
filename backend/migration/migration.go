package migrate_db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/ireuven89/hello-world/backend/environment"
)

func Migrate() error {
	url := environment.Variables.DbUrl
	service, err := migrate.New("/migrations/", url)

	service.m

	if err != nil {
		return err
	}

	return nil
}
