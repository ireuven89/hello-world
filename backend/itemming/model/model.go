package model

import "time"

type Item struct {
	ID          string      `json:"id" db:"id"`
	UserUuid    string      `json:"userUuid" db:"user_uuid"`
	Link        interface{} `json:"link" db:"link"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Item        string      `json:"itemming" db:"itemming"`
	Price       int64       `json:"price" db:"price"`
	CreatedAt   time.Time   `json:"CreatedAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"UpdatedAt" db:"updated_at"`
}

type Category int

const (
	Art Category = iota
)

type ItemInput struct {
	Uuid        string `json:"uuid"`
	UserUuid    string `json:"userUuid"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
}

type ListInput struct {
	Price       string
	Name        string
	Description string
}
