package models

import (
	"time"

	"bitbucket.org/kiloops/api/gitadmin"
	"github.com/jinzhu/gorm"
)

const (
	_ = iota
	Creator
	Collaborator
)

//UsersProjects ManyToMany rel
type UsersProjects struct {
	ID        int     `gorm:"primary_key" json:"id"`
	Role      int     `json:"role"`
	ProjectID int     `json:"project_id"`
	UserID    int     `json:"user_id"`
	Project   Project `json:"project"`
	User      Project `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName for UsersProjects
func (u UsersProjects) TableName() string {
	return "users_projects"
}

// AfterCreate callback
func (u UsersProjects) AfterCreate(txn *gorm.DB) error {
	return u.withRels(txn, func(email string, SSHs *[]SSH, slug string) error {
		path := gitadmin.ProjectPath(slug)
		for _, ssh := range *SSHs {
			if err := gitadmin.AddSSHToProject(email, ssh.ID, slug); err != nil {
				gitadmin.RevertChange(path)
				return err
			}
			gitadmin.CommitChange(path)
		}
		return nil
	})
}

// AfterDelete callback
func (u UsersProjects) AfterDelete(txn *gorm.DB) error {
	err := u.withRels(txn, func(email string, SSHs *[]SSH, slug string) error {
		path := gitadmin.ProjectPath(slug)
		for _, ssh := range *SSHs {
			if err := gitadmin.RemoveSSHFromProject(email, ssh.ID, slug); err != nil {
				gitadmin.RevertChange(path)
				return err
			}
		}
		gitadmin.CommitChange(path)
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
		return ErrUserNotFound
	}
	if project.ID == 0 {
		return ErrProjectNotFound
	}
	var SSHs []SSH
	txn.Model(&user).Related(&SSHs)
	return f(user.Email, &SSHs, project.Slug)
}
