package models

import "time"

type Transaction struct {
	ID              string    `json:"id" gorm:"default:gen_random_uuid()"`
	AccountNumber   string    `json:"account_number"`
	ActionPerformed string    `json:"action_performed"`
	Recipient       string    `json:"recipient"`
	Balance         float64   `json:"balance"`
	CreatedAt       time.Time `json:"created_at" sql:"type:timestamp without time zone"`
	CreatedDate     string    `json:"created_date"`
}
