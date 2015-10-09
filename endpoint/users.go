package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/kiloops/api/ierrors"
	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//UserCreate serves the route POST /users
func UserCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var userLogin models.UserLogin
		if err := c.BindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return false
		}
		if valid, errMap := models.ValidStruct(&userLogin); !valid {
			errorResponseMap(c, errMap)
			return false
		}
		user := models.User{Email: userLogin.Email, Password: userLogin.Password}
		if err := txn.Create(&user).Error; err != nil {
			errorResponse(c, err, ierrors.ErrUserCreate)
			return false
		}
		userLogged := models.UserLogged{ID: user.ID, Email: user.Email, Token: user.Token}
		c.JSON(http.StatusOK, userLogged)
		return true
	})
}

//UserLogin serves the route POST /users/login
func UserLogin(c *gin.Context) {
	var userLogin models.UserLogin
	if err := c.BindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var user models.User
	models.Gdb.Find(&user, "email = ?", userLogin.Email)
	switch {
	case user.Email == "":
		c.JSON(http.StatusNotFound, "")
	case !user.LoggedIn(userLogin):
		errorResponse(c, ierrors.ErrUserLogin)
	case !user.Active:
		errorResponse(c, ierrors.ErrUserInactive)
	default:
		userLogged := models.UserLogged{ID: user.ID, Email: user.Email, Token: user.Token, Active: user.Active}
		c.JSON(http.StatusOK, userLogged)
	}
}

//UserForgot serves the route POST /users/forgot
func UserForgot(c *gin.Context) {
	var email models.Email
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if valid, errMap := models.ValidStruct(&email); !valid {
		errorResponseMap(c, errMap)
		return
	}
	models.InTx(func(txn *gorm.DB) bool {
		var user models.User
		models.Gdb.Find(&user, "email like ?", email.Value)
		if user.ID == 0 {
			c.JSON(http.StatusNotFound, "User not found")
			return false
		}
		passwordRequest := models.PasswordRequest{}
		if err := txn.Create(&passwordRequest).Error; err != nil {
			fmt.Println(err)
			errorResponse(c, err, ierrors.ErrPasswordRequestCreate)
			return false
		}
		//send email
		go utils.SendChangePasswordEmail(user.Email, user.Email, passwordRequest.Token)
		c.JSON(http.StatusOK, passwordRequest)
		return true
	})
}

//UserChangePassword serves the route POST /users/change_password
func UserChangePassword(c *gin.Context) {
	var changePassword models.ChangePassword
	if err := c.BindJSON(&changePassword); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if valid, errMap := models.ValidStruct(&changePassword); !valid {
		errorResponseMap(c, errMap)
		return
	}
	models.InTx(func(txn *gorm.DB) bool {
		var user models.User
		models.Gdb.Find(&user, "email like ?", changePassword.Email)
		if user.ID == 0 {
			c.JSON(http.StatusNotFound, "User not found")
			return false
		}

		var passwordRequest models.PasswordRequest
		models.Gdb.Find(&passwordRequest, "token like ? and used = false", changePassword.Token)
		if passwordRequest.ID == 0 {
			c.JSON(http.StatusNotFound, "Password Request not found")
			return false
		}

		passwordRequest.Used = true
		if err := txn.Save(&passwordRequest).Error; err != nil {
			errorResponse(c, err, ierrors.ErrPasswordRequestUpdate)
			return false
		}

		//validate time
		duration := time.Since(passwordRequest.CreatedAt)
		if duration.Hours() > 1.0 {
			errorResponse(c)
		}

		user.SetPassword(changePassword.Password)
		if err := txn.Save(&user).Error; err != nil {
			errorResponse(c, err, ierrors.ErrUserUpdate)
			return false
		}

		c.JSON(http.StatusOK, "")
		return false
	})
}
