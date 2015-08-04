package models

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
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
	u.Password = MD5(u.Password)
	u.Token = GenerateToken(20)
}

//LoggedIn validtes if a user is logged
func (u User) LoggedIn(login UserLogin) bool {
	return MD5(login.Password) == u.Password
}

//MD5 encription
func MD5(cad string) string {
	hash := sha1.New()
	hash.Write([]byte(cad))
	return hex.EncodeToString(hash.Sum(nil))
}

//GenerateToken a random
func GenerateToken(size int) string {
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(rb)
}
