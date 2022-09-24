package models

import "time"

// store transactions with currency or add field for currency
type Transaction struct {
	ID              string    `json:"id" gorm:"default:gen_random_uuid()"`
	AccountNumber   string    `json:"account_number"`
	ActionPerformed string    `json:"action_performed"`
	Recipient       string    `json:"recipient,omitempty"`
	Balance         float64   `json:"balance"`
	CreatedAt       time.Time `json:"created_at" sql:"type:timestamp without time zone"`
	CreatedDate     string    `json:"created_date"`
}
