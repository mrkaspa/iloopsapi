package models

import "time"

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}
