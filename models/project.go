package models

import "time"

type Project struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
