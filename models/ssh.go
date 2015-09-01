package models

import (
	"time"

	"bitbucket.org/kiloops/api/gitadmin"

	"github.com/jinzhu/gorm"
	"github.com/mrkaspa/go-helpers"
)

//SSH key
type SSH struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Name      string    `json:"name" validate:"required"`
	PublicKey string    `sql:"type:text" json:"public_key" validate:"required"`
	Hash      string    `sql:"type:varchar(100);unique" json:"-"`
	UserID    int       `json:"user_id"`
	User      User      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

//BeforeCreate callback
func (s *SSH) BeforeCreate() error {
	s.Hash = helpers.MD5(s.PublicKey)
	return nil
}

//AfterCreate callback
func (s *SSH) AfterCreate(txn *gorm.DB) error {
	var user User
	txn.Model(s).Related(&user)
	if user.ID == 0 {
		return ErrUserNotFound
	}
	if err := gitadmin.AddSSH(user.Email, s.ID, s.PublicKey); err != nil {
		return err
	}
	userProjects := user.AllProjects()
	for _, userProject := range *userProjects {
		err := gitadmin.InTx(func() error { return gitadmin.AddSSHToProject(user.Email, s.ID, userProject.Project.Slug) })
		if err != nil {
			return err
		}
	}
	return nil
}

//AfterDelete callback
func (s *SSH) AfterDelete(txn *gorm.DB) error {
	var user User
	txn.Model(s).Related(&user)
	if user.ID == 0 {
		return ErrUserNotFound
	}
	if err := gitadmin.DeleteSSH(user.Email, s.ID); err != nil {
		return err
	}
	userProjects := user.AllProjects()
	for _, userProject := range *userProjects {
		err := gitadmin.InTx(func() error { return gitadmin.RemoveSSHFromProject(user.Email, s.ID, userProject.Project.Slug) })
		if err != nil {
			return err
		}
	}
	return nil
}

//TableName for SSH
func (s SSH) TableName() string {
	return "ssh"
}
