package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/kiloops/api/utils"
)

//User model
type User struct {
	ID        int `gorm:"primary_key"`
	Email     string
	Password  string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//UserLogin model
type UserLogin struct {
	Email    string `json:"email" validate:"nonzero,regexp=^[A-Za-z0-9._%+-]+@[A-Z0-9.-]+\.[A-Za-z]{2,4}$"`
	Password string `json:"password" validate:"nonzero"`
}

//UserLogged model
type UserLogged struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

//BeforeCreate callback
func (u *User) BeforeCreate() {
	u.Password = utils.MD5(u.Password)
	u.Token = utils.GenerateToken(20)
}

//LoggedIn validtes if a user is logged
func (u User) LoggedIn(login UserLogin) bool {
	return utils.MD5(login.Password) == u.Password
}

func (u User) AllProjects() []UsersProjects {
	userProjects := []UsersProjects{}
	Gdb.Where("user_id = ?", u.ID).Find(&userProjects)
	return userProjects
}

func (u User) OwnedProjects() []UsersProjects {
	userProjects := []UsersProjects{}
	Gdb.Where("user_id = ? and role = ?", u.ID, Creator).Find(&userProjects)
	return userProjects
}

func (u User) CollaboratorProjects() []UsersProjects {
	userProjects := []UsersProjects{}
	Gdb.Where("user_id = ? and role = ?", u.ID, Collaborator).Find(&userProjects)
	return userProjects
}

func (u User) CreateProject(txn *gorm.DB, project *Project) error {
	if txn.Save(&project).Error == nil {
		//Creates a relation between the user and the project
		userProject := UsersProjects{Role: Creator, UserID: u.ID, ProjectID: project.ID}
		if txn.Save(&userProject).Error == nil {
			return nil
		} else {
			txn.Rollback()
			return errors.New("User Project can't be saved")
		}
	} else {
		return errors.New("Project can't be saved")
	}
}

func (u User) HasAdminAccessTo(projectID int) bool {
	var count int
	Gdb.Where("user_id = ? and project_id = ? and role = ?", u.ID, projectID, Creator).Count(&count)
	return count > 0
}

func (u User) HasWriteAccessTo(projectID int) bool {
	var count int
	Gdb.Where("user_id = ? and project_id = ? and role in (?,?)", u.ID, projectID, Creator, Collaborator).Count(&count)
	return count > 0
}
