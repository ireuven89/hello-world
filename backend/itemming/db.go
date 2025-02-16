package itemming

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/ido50/sqlz"

	"github.com/ireuven89/hello-world/backend/environment"
)

func MustNewDB() (*sqlz.DB, string, error) {
	cfg := mysql.Config{
		User:   environment.Variables.ItemsDbUser,
		Passwd: environment.Variables.ItemsDbPassword,
		Addr:   environment.Variables.ItemsDbHost,
		DBName: "items",
		Net:    "tcp",
	}
	itemsDB, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Printf("failed to open db items: %v", err)
		return nil, "", err
	}

	//ping check
	if err = itemsDB.Ping(); err != nil {
		fmt.Printf("failed to dial to db items: %v", err)
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

	return sqlz.New(itemsDB, "sql"), migrationPath, nil
}
