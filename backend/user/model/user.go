package user

import (
	"github.com/ireuven89/hello-world/backend/routes"
)

type User struct {
	ID     int    `json:"ID" db:"id"`
	Name   string `json:"Name" db:"name"`
	Region string `json:"Region" db:"region"`
}

type UserRequest struct {
	routes.Pagination
	Name   string `json:"Name"`
	Region string `json:"Region"`
}
