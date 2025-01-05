package migration

import (
	"github.com/golang-migrate/migrate/v4"
)

func Migrate(url string) error {
	service, err := migrate.New("/migrations/", url)

	err = service.Up()

	if err != nil {
		return err
	}

	return nil
}
