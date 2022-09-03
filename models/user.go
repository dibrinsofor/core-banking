package models

import "time"

type User struct {
	AccountNumber string    `json:"account_number" gorm:"default:gen_random_uuid()"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Balance       int       `json:"balance"`
	CreatedAt     time.Time `json:"created_at" sql:"type:timestamp without time zone"`
}
