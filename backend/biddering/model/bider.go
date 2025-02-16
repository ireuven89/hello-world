package model

import (
	"strconv"
	"time"
)

type Bidder struct {
	Id          int64     `json:"-" db:"id"`
	Uuid        string    `json:"uuid" db:"uuid"`
	UserUuid    string    `json:"UserUuid" db:"user_uuid"`
	Name        string    `json:"Name" db:"name"`
	Item        string    `json:"Item" db:"itemming"`
	Price       int64     `json:"Price" db:"price"`
	Description string    `json:"Description" db:"description"`
	CreatedAt   time.Time `json:"Created_At" db:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" db:"updated_at"`
}

type BiddersInput struct {
	Page PageRequest `json:"defaultRequest"`
	Uuid string      `json:"uuid"`
	Name string      `json:"name"`
	Item string      `json:"itemming"`
}

type BidderInput struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Item        string `json:"itemming"`
	UserUuid    string `json:"UserUuid"`
	Price       int64  `json:"Price"`
	Description string `json:"Description"`
}

type PageRequest struct {
	Offset int64
	Limit  int64
}

func SetPage(offset string, limit string) PageRequest {
	parsedOffset, err := strconv.Atoi(offset)

	parsedLimit, err := strconv.Atoi(limit)

	if err != nil {
		parsedLimit = 50
	}

	return PageRequest{
		Offset: int64(parsedOffset),
		Limit:  int64(parsedLimit),
	}
}

func (p *PageRequest) GetLimit() int64 {
	if p.Limit == 0 {
		return 50
	}

	return p.Limit
}
