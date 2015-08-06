package models

import (
	"time"

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
