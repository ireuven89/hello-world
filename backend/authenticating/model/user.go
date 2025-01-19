package model

type User struct {
	Username string `db:"user"`
	Password string `db:"password"` // This will store the hashed password
}

const TableName = "users"
