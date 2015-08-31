package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/models"
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
			c.JSON(http.StatusConflict, errMap)
			return false
		}
		user := models.User{Email: userLogin.Email, Password: userLogin.Password}
		if txn.Create(&user).Error != nil {
			c.JSON(http.StatusBadRequest, "User can't be saved")
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
		c.JSON(http.StatusBadRequest, "Could not authenticate the User")
		return
	}
	var user models.User
	err := models.Gdb.Find(&user, "email = ?", userLogin.Email).Error
	if err == nil && user.Email != "" && user.LoggedIn(userLogin) {
		userLogged := models.UserLogged{ID: user.ID, Email: user.Email, Token: user.Token}
		c.JSON(http.StatusOK, userLogged)
	} else {
		c.JSON(http.StatusBadRequest, "Could not authenticate the User")
	}
}
