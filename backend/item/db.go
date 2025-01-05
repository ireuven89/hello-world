package item

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/ireuven89/hello-world/backend/environment"
	"path/filepath"
)

func MustNewDB() (*sql.DB, string, error) {
	cfg := mysql.Config{
		User:   environment.Variables.ItemsDbUser,
		Passwd: environment.Variables.ItemsDbPassword,
		Addr:   environment.Variables.ItemsDbHost,
		DBName: "items",
		Net:    "tcp",
	}
	itemsDB, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, "", err
	}

	//ping check
	if err = itemsDB.Ping(); err != nil {
		return nil, "", err
	}

	//create lock table if not exists
	if _, err = itemsDB.Exec("create table if not exists lock_table(lock_row int)"); err != nil {
		return nil, "", err
	}

	//set migration dir
	migrationPath, err := filepath.Abs("./db/migrations/items")

	if err != nil {
		return nil, "", err
	}

	return itemsDB, migrationPath, nil
}
