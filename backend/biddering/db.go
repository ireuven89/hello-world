package biddering

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
		DBName: "bidders",
		Net:    "tcp",
	}
	add := cfg.FormatDSN()
	bidderDB, err := sql.Open("mysql", add)

	if err != nil {
		fmt.Printf("bidderring.MustNewDB failed to dial to db bidders: %v", err)
		return nil, "", err
	}

	//ping check
	if err = bidderDB.Ping(); err != nil {
		fmt.Printf("bidderring.MustNewDB failed to dial to db bidders: %v", err)
		return nil, "", err
	}

	//create table if not exists
	if _, err = bidderDB.Exec("CREATE TABLE IF NOT EXISTS lock_table (lock_row int)"); err != nil {
		fmt.Printf("bidderring.MustNewDB failed to create lock table: %v", err)
		return nil, "", err
	}

	//set migration dir
	migrationDir, err := filepath.Abs("./db/migrations/bidders")

	if err != nil {
		return nil, "", err
	}

	return sqlz.New(bidderDB, "mysql"), migrationDir, nil
}
