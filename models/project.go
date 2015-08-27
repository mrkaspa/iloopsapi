package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/utils"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

//Project on the system
type Project struct {
	ID          int    `gorm:"primary_key" json:"id"`
	Slug        string `json:"slug" sql:"unique"`
	Name        string `json:"name" validate:"required"`
	URLRepo     string `json:"url_repo"`
	Periodicity string `json:"periodicity"`
	Command     string `json:"command"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//AfterCreate a Project
func (p *Project) AfterCreate(txn *gorm.DB) error {
	p.SetSlug()
	p.URLRepo = gitadmin.GITURLROOT + ":" + p.Slug + ".git"
	if err := txn.Save(p).Error; err == nil {
		return gitadmin.CreateProject(p.Slug)
	} else {
		return err
	}
}

//SetSlug for the project
func (p *Project) SetSlug() {
	nameSlug := slug.Make(p.Name)
	p.Slug = fmt.Sprintf("%s-%d", nameSlug, p.ID)
}

//BeforeDelete a Project
func (p *Project) BeforeDelete(txn *gorm.DB) error {
	return txn.Where("project_id = ?", p.ID).Delete(UsersProjects{}).Error
	// var userProjects []UsersProjects
	// txn.Where("project_id = ?", p.ID).Find(&userProjects)
	// for _, userProject := range userProjects {
	// 	if err := txn.Delete(&userProject).Error; err != nil {
	// 		return err
	// 	}
	// }
	// return nil
}

//AfterDelete a Project
func (p *Project) AfterDelete() error {
	return gitadmin.DeleteProject(p.Slug)
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
	return ErrCreatorNotRemoved
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
		return ErrUserIsNotCollaborator
	}
}

func (p Project) GetCommand() string {
	return fmt.Sprintf("docker run -d %s", p.Slug)
}

func (p Project) Schedule() error {
	task := Task{ID: p.Slug, Periodicity: p.Periodicity, Command: p.Command}
	taskJSON, _ := json.Marshal(task)
	client := utils.Client{
		&http.Client{},
		"http://" + os.Getenv("GUARTZ_HOST"),
		"application/json",
	}
	resp, err := client.CallRequest("POST", "/tasks", bytes.NewReader(taskJSON))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrTaskNotScheduled
	}
	return nil
}

//FindProject by id
func FindProject(id int) (*Project, error) {
	var project Project
	Gdb.First(&project, id)
	if project.ID != 0 {
		return &project, nil
	}
	return nil, ErrProjectNotFound
}

func FindProjectBySlug(slug string) (*Project, error) {
	var project Project
	Gdb.Where("slug like ?", slug).First(&project)
	if project.ID != 0 {
		return &project, nil
	}
	return nil, ErrProjectNotFound
}
