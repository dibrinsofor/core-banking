package models

import "time"

type AccountInfo struct {
	AccountNumber  string    `json:"account_number" gorm:"default:gen_random_uuid()"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Balance        float64   `json:"balance,omitempty"`
	AccountBalance string    `json:"account_balance"`
	CreatedAt      time.Time `json:"created_at" sql:"type:timestamp without time zone"`
}
