package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Project struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Project) DeleteRels(txn *gorm.DB) {
	txn.Where("project_id = ?", p.ID).Delete(UsersProjects{})
}
