package models

import (
	"errors"
	"time"
)

//Project on the system
type Project struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//DeleteRels UsersProjects
func (p *Project) DeleteRels(txn *KDB) {
	txn.Where("project_id = ?", p.ID).Delete(UsersProjects{})
}

//AddUser adds new user
func (p *Project) AddUser(txn *KDB, user *User) error {
	r := UsersProjects{Role: Collaborator, UserID: user.ID, ProjectID: p.ID}
	return txn.Save(&r).Error
}

//DelegateUser sets an user as Creator
func (p *Project) DelegateUser(txn *KDB, userAdmin, user *User) error {
	if user.HasCollaboratorAccessTo(p.ID) {
		if err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", userAdmin.ID, p.ID).Update("role", Collaborator).Error; err == nil {
			if err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", user.ID, p.ID).Update("role", Creator).Error; err == nil {
				return nil
			} else {
				txn.KRollback()
				return err
			}
		} else {
			txn.KRollback()
			return err
		}
	} else {
		txn.KRollback()
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
