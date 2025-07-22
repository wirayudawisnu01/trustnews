package entity

import "time"

type ContentEntity struct {
	ID          int64
	Title       string
	Excerpt     string
	Description string
	Image       string
	Tags        []string
	Status      string
	IsValid     string
	CategoryID  int64
	CreatedByID int64
	CreatedAt   time.Time
	Category 	CategoryEntity
	User 		UserEntity
}

type QueryString struct {
	Limit 		int
	Page 		int
	OrderBy 	string
	OrderType	string
	Search 		string
	CategoryID	int64
	Status		string
}