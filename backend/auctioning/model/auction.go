package model

import "time"

type Auction struct {
	ID               string    `json:"-" db:"id"`
	Uuid             string    `json:"uuid" db:"uuid"`
	Item             string    `json:"item" db:"item"`
	Price            int64     `json:"price" db:"price"`
	WinningPrice     int64     `json:"winningPrice" db:"winning_price"`
	UserUuid         string    `json:"UserUuid" db:"user_uuid"`
	BiddersCount     int64     `json:"biddersCount" db:"bidders_count"`
	BiddersThreshold int64     `json:"biddersThreshold" db:"bidders_threshold"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	ExpiredAt        time.Time `json:"expiredAt" db:"expired_at"`
	Status           Status    `json:"status" db:"status"`
}

type AuctionRequest struct {
	Id               string `json:"id"`
	Category         string `json:"category"`
	Price            int64  `json:"price" db:"price"`
	WinningPrice     int64  `json:"winningPrice" db:"winning_price"`
	UserUuid         string `json:"UserUuid" db:"user_uuid"`
	BiddersCount     int64  `json:"biddersCount" db:"bidders_count"`
	BiddersThreshold int64  `json:"biddersThreshold" db:"bidders_threshold"`
	Status           string `json:"status" db:"status"`
}

type Status int

const (
	InProgress Status = iota
	Sold
	Expired
)
