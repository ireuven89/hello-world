package authenticating

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/ido50/sqlz"

	"github.com/ireuven89/hello-world/backend/utils"
)

func MustNewDB(config utils.DataBaseConnection) (*sqlz.DB, string, error) {
	cfg := mysql.Config{
		User:   os.Getenv(config.UserName),
		Passwd: os.Getenv(config.Password),
		Addr:   os.Getenv(config.Host),
		DBName: "auth",
		Net:    "tcp",
	}
	add := cfg.FormatDSN()

	fmt.Printf(add)
	authDB, err := sql.Open("mysql", add)

	if err != nil {
		fmt.Printf("authenticating.MustNewDB failed to dial to db auth: host %s address: %v", os.Getenv(config.Host), err)
		return nil, "", err
	}

	//ping check
	if err = authDB.Ping(); err != nil {
		fmt.Printf("authenticating.MustNewDB failed to open db host: %s, %v", os.Getenv(config.Host), err)
		return nil, "", err
	}

	//create lock table if not exists
	if _, err = authDB.Exec("create table if not exists lock_table(lock_row int)"); err != nil {
		return nil, "", err
	}

	//set migration dir
	migrationDir, err := filepath.Abs("./db/migrations/auth")

	if err != nil {
		return nil, "", err
	}

	return sqlz.New(authDB, "mysql"), migrationDir, nil
}
