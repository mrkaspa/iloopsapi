package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mrkaspa/go-helpers"
)

type PasswordRequest struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Token     string    `json:"token"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

//BeforeCreate callback
func (p *PasswordRequest) AfterCreate(txn *gorm.DB) error {
	rand := strconv.Itoa(p.ID) + helpers.RandomString(20)
	p.Token = helpers.MD5(rand)
	if err := txn.Save(p).Error; err != nil {
		return err
	}
	return nil
}
