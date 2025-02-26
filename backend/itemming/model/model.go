package model

type Item struct {
	ID          int64       `json:"-" db:"id"`
	Uuid        string      `json:"uuid" db:"uuid"`
	UserUuid    string      `json:"userUuid" db:"user_uuid"`
	Link        interface{} `json:"link" db:"link"`
	Category    string      `json:"category" db:"category"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
}

type Category int

const (
	Art Category = iota
)

type ItemInput struct {
	Uuid     string `json:"uuid" db:"uuid"`
	UserUuid string `json:"userUuid" db:"user_uuid"`
	Category string `json:"category" db:"category"`
	Name     string `json:"name" db:"name"`
}

type ListInput struct {
	Link        string
	Name        string
	Description string
	Page        int
	Size        int
}
