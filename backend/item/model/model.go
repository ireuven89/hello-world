package item

type Item struct {
	ID       int64       `json:"-" db:"id"`
	Uuid     string      `json:"uuid" db:"uuid"`
	UserUuid string      `json:"userUuid" db:"user_uuid"`
	Link     interface{} `json:"link" db:"link"`
	Category string      `json:"category" db:"category"`
}

type Category int

const (
	Art Category = iota
)
