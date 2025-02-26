package userring

import (
	"database/sql"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/ido50/sqlz"

	"github.com/ireuven89/hello-world/backend/environment"
)

// MustNewDB - returns the db connection, the migrations directory of the db, and an error if anything failed
func MustNewDB() (*sqlz.DB, string, error) {
	cfg := mysql.Config{
		User:   environment.Variables.UsersDbUser,
		Passwd: environment.Variables.UsersDbPassword,
		Addr:   environment.Variables.UsersDbHost,
		DBName: "userring",
		Net:    "tcp",
	}
	add := cfg.FormatDSN()
	usersDB, err := sql.Open("mysql", add)
	if err != nil {
		return nil, "", err
	}

	//ping check
	if err = usersDB.Ping(); err != nil {
		return nil, "", err
	}

	//create lock table if not exists
	if _, err = usersDB.Exec("create table if not exists lock_table(lock_row int)"); err != nil {
		return nil, "", err
	}

	//set migration dir
	migrationDir, err := filepath.Abs("./db/migrations/userring")

	if err != nil {
		return nil, "", err
	}

	return sqlz.New(usersDB, "mysql"), migrationDir, nil
}
