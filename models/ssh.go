package models

import (
	"time"

	"bitbucket.org/kiloops/api/utils"
)

//SSH key
type SSH struct {
	ID        int       `gorm:"primary_key" json:"id"`
	PublicKey string    `sql:"type:text" json:"public_key" validate:"required"`
	Hash      string    `sql:"type:varchar(500)" json:"-"`
	UserID    int       `json:"user_id"`
	User      User      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

//BeforeCreate callback
func (s *SSH) BeforeCreate() error {
	s.Hash = utils.MD5(s.PublicKey)
	return nil
}

//TableName for SSH
func (s SSH) TableName() string {
	return "ssh"
}
