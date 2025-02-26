package model

type User struct {
	ID       string `db:"id"`
	Username string `db:"user"`
	Password string `db:"password"` // This will store the hashed password
}

type Page struct {
	Page     int64
	PageSize int64
}

const TableName = "userring"
