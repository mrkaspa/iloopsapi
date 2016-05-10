package models

import (
	"time"

	"github.com/mrkaspa/iloopsapi/gitadmin"
	"github.com/mrkaspa/iloopsapi/ierrors"
	"github.com/jinzhu/gorm"
)

const (
	_ = iota
	Creator
	Collaborator
)

const maxNumProjects = 5

//UsersProjects ManyToMany rel
type UsersProjects struct {
	ID        int     `gorm:"primary_key" json:"id"`
	Role      int     `json:"role" validate:"required"`
	ProjectID int     `json:"project_id" validate:"required"`
	UserID    int     `json:"user_id" validate:"required"`
	Project   Project `json:"project"`
	User      Project `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName for UsersProjects
func (u UsersProjects) TableName() string {
	return "users_projects"
}

func (u *UsersProjects) BeforeCreate(txn *gorm.DB) error {
	var counter int
	txn.Model(UsersProjects{}).Where("user_id = ?", u.UserID).Count(&counter)
	if counter >= maxNumProjects {
		return ierrors.ErrUserExceedMaxProjects
	}
	return nil
}

// AfterCreate callback
func (u *UsersProjects) AfterCreate(txn *gorm.DB) error {
	return u.withRels(txn, func(email string, SSHs *[]SSH, slug string) error {
		for _, ssh := range *SSHs {
			err := gitadmin.InTx(func() error { return gitadmin.AddSSHToProject(email, ssh.ID, slug) })
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// AfterDelete callback
func (u *UsersProjects) AfterDelete(txn *gorm.DB) error {
	err := u.withRels(txn, func(email string, SSHs *[]SSH, slug string) error {
		for _, ssh := range *SSHs {
			err := gitadmin.InTx(func() error { return gitadmin.RemoveSSHFromProject(email, ssh.ID, slug) })
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (u UsersProjects) withRels(txn *gorm.DB, f func(string, *[]SSH, string) error) error {
	var project Project
	var user User
	txn.Model(&u).Related(&user)
	txn.Model(&u).Related(&project)
	if user.ID == 0 {
		return ierrors.ErrUserNotFound
	}
	if project.ID == 0 {
		return ierrors.ErrProjectNotFound
	}
	var SSHs []SSH
	txn.Model(&user).Related(&SSHs)
	return f(user.Email, &SSHs, project.Slug)
}
