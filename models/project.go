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

var guartzClient utils.Client

//InitGuartzClient with os env
func InitGuartzClient() {
	guartzClient = utils.Client{
		&http.Client{},
		"http://" + os.Getenv("GUARTZ_HOST"),
		"application/json",
	}
}

//AfterCreate a Project
func (p *Project) AfterCreate(txn *gorm.DB) error {
	p.SetSlug()
	p.URLRepo = gitadmin.GITURLROOT + ":" + p.Slug + ".git"
	if err := txn.Save(p).Error; err != nil {
		return err
	}
	return gitadmin.CreateProject(p.Slug)
}

//SetSlug for the project
func (p *Project) SetSlug() {
	nameSlug := slug.Make(p.Name)
	p.Slug = fmt.Sprintf("%s-%d", nameSlug, p.ID)
}

//BeforeDelete a Project
func (p *Project) BeforeDelete(txn *gorm.DB) error {
	go p.Stop()
	return txn.Where("project_id = ?", p.ID).Delete(UsersProjects{}).Error
}

//AfterDelete a Project
func (p *Project) AfterDelete() error {
	return gitadmin.DeleteProject(p.Slug)
}

//AddUser adds new user
func (p *Project) AddUser(txn *gorm.DB, user *User) error {
	r := UsersProjects{Role: Collaborator, UserID: user.ID, ProjectID: p.ID}
	if err := txn.Create(&r).Error; err != nil {
		if err != ErrUserExceedMaxProjects {
			return ErrUserProjectNotSaved
		}
		return err
	}
	return nil
}

//RemoveUser removes and user
func (p *Project) RemoveUser(txn *gorm.DB, user *User) error {
	var userProject UsersProjects
	txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", user.ID, p.ID).First(&userProject)
	if userProject.Role != Collaborator {
		return ErrCreatorNotRemoved
	}
	return txn.Delete(&userProject).Error
}

//DelegateUser sets an user as Creator
func (p *Project) DelegateUser(txn *gorm.DB, userAdmin, user *User) error {
	if !user.HasCollaboratorAccessTo(p.ID) {
		return ErrUserIsNotCollaborator
	}
	if err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", userAdmin.ID, p.ID).Update("role", Collaborator).Error; err != nil {
		return err
	}
	err := txn.Model(UsersProjects{}).Where("user_id = ? and project_id = ?", user.ID, p.ID).Update("role", Creator).Error
	if err != nil {
		return err
	}
	return nil
}

//GetCommand create the command
func (p Project) GetCommand() string {
	return fmt.Sprintf("docker run -d %s:latest", p.Slug)
}

//Schedule a worker
func (p Project) Schedule() error {
	task := Task{ID: p.Slug, Periodicity: p.Periodicity, Command: p.Command}
	taskJSON, _ := json.Marshal(task)
	return guartzClient.CallRequest("POST", "/tasks", bytes.NewReader(taskJSON)).WithResponse(func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return ErrTaskNotScheduled
		}
		return nil
	})
}

//Stop a worker
func (p *Project) Stop() error {
	return guartzClient.CallRequestNoBody("DELETE", "/tasks/"+p.Slug).WithResponse(func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			return ErrTaskNotStopped
		}
		return nil
	})
}

//FindProject by id
func FindProject(id int) (*Project, error) {
	var project Project
	Gdb.First(&project, id)
	if project.ID == 0 {
		return nil, ErrProjectNotFound
	}
	return &project, nil
}

//FindProjectBySlug by slug
func FindProjectBySlug(slug string) (*Project, error) {
	var project Project
	Gdb.Where("slug like ?", slug).First(&project)
	if project.ID == 0 {
		return nil, ErrProjectNotFound
	}
	return &project, nil
}
