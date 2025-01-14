package model

import (
	"github.com/ireuven89/hello-world/backend/routes"
)

type User struct {
	ID   int    `json:"-" db:"id"`
	Uuid string `json:"uuid" db:"uuid"`
	Name string `json:"name" db:"name"`
}

func (u *User) IsEmpty() bool {

	return u.ID == 0 && u.Uuid == "" && u.Name == ""
}

type UserFetchInput struct {
	routes.Pagination
	Name   string `json:"name"`
	Uuid   string `json:"uuid"`
	Region string `json:"region"`
}

type UserUpsertInput struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Email  string `json:"email"`
}

type DeleteUserInput struct {
	Uuid string `json:"uuid"`
}

type UserResponse struct {
	Name        string `json:"name"`
	Uuid        string `json:"uuid"`
	Region      string `json:"region"`
	Description string `json:"description"`
}
