package biddering

import (
	"database/sql"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/ido50/sqlz"

	"github.com/ireuven89/hello-world/backend/utils"
)

func MustNewDB(config utils.DataBaseConnection) (*sqlz.DB, string, error) {
	cfg := mysql.Config{
		User:   config.UserName,
		Passwd: config.Password,
		Addr:   config.Host,
		DBName: "bidder",
		Net:    "tcp",
	}
	add := cfg.FormatDSN()
	bidderDB, err := sql.Open("mysql", add)

	if err != nil {
		return nil, "", err
	}

	//ping check
	if err = bidderDB.Ping(); err != nil {
		return nil, "", err
	}

	//create table if not exists
	if _, err = bidderDB.Exec("CREATE TABLE IF NOT EXISTS lock_table (lock_row int)"); err != nil {
		return nil, "", err
	}

	//set migration dir
	migrationDir, err := filepath.Abs("./db/migrations/bidders")

	if err != nil {
		return nil, "", err
	}

	return sqlz.New(bidderDB, "mysql"), migrationDir, nil
}
