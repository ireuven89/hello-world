package model

import "github.com/ireuven89/hello-world/backend/utils"

type Bidder struct {
	Id          int64  `json:"-" db:"id"`
	Uuid        string `json:"uuid" db:"uuid"`
	UserUuid    string `json:"UserUuid" db:"user_uuid"`
	Name        string `json:"Name" db:"name"`
	Item        string `json:"Item" db:"item"`
	Price       string `json:"Price" db:"price"`
	Description string `json:"Description" db:"description"`
}

type BiddersInput struct {
	DefaultRequest utils.DefaultRequest `json:"defaultRequest"`
	Uuid           string               `json:"uuid" db:"uuid"`
	Name           string               `json:"name"`
	Item           string               `json:"item"`
}
