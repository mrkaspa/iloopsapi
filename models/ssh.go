package models

import (
	"time"

	"bitbucket.org/kiloops/api/utils"
)

//SSH key
type SSH struct {
	ID        int       `gorm:"primary_key" json:"id"`
	PublicKey string    `sql:"type:text" json:"public_key" validate:"nonzero"`
	Hash      string    `sql:"type:varchar(500)" json:"-"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

//BeforeCreate callback
func (s *SSH) BeforeCreate() {
	s.Hash = utils.MD5(s.PublicKey)
}

//TableName for SSH
func (s SSH) TableName() string {
	return "ssh"
}