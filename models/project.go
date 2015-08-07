package models

import "time"

type Project struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Token  string `sql:"-" json:"-"`
	UserID int    `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
