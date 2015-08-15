package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

//Project on the system
type Project struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//AfterCreate a Project
func (p *Project) AfterCreate(txn *gorm.DB) error {
	nameSlug := slug.Make(p.Name)
	p.Slug = fmt.Sprintf("%s-%d", nameSlug, p.ID)
	return txn.Save(p).Error
}

//BeforeDelete a Project
func (p *Project) BeforeDelete(txn *gorm.DB) error {
	return txn.Where("project_id = ?", p.ID).Delete(UsersProjects{}).Error
}

//AddUser adds new user
func (p *Project) AddUser(txn *gorm.DB, user *User) error {
	r := UsersProjects{Role: Collaborator, UserID: user.ID, ProjectID: p.ID}
	return txn.Save(&r).Error
}

//RemoveUser removes and user
func (p *Project) RemoveUser(txn *gorm.DB, user *User) error {
	var userProject UsersProjects
	txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", user.ID, p.ID).First(&userProject)
	if userProject.Role == Collaborator {
		return txn.Delete(&userProject).Error
	}
	return errors.New("You can't remove a Creator from a project")
}

//DelegateUser sets an user as Creator
func (p *Project) DelegateUser(txn *gorm.DB, userAdmin, user *User) error {
	if user.HasCollaboratorAccessTo(p.ID) {
		if err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", userAdmin.ID, p.ID).Update("role", Collaborator).Error; err == nil {
			if err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", user.ID, p.ID).Update("role", Creator).Error; err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("The user doesn't have collaborator access to the project")
	}
}

//FindProject by id
func FindProject(id int) (*Project, error) {
	var project Project
	Gdb.First(&project, id)
	if project.ID != 0 {
		return &project, nil
	}
	return nil, errors.New("Project not found")
}
