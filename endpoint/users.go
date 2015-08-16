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
		if err := c.BindJSON(&userLogin); err == nil {
			if valid, errMap := models.ValidStruct(&userLogin); valid {
				user := models.User{Email: userLogin.Email, Password: userLogin.Password}
				if txn.Save(&user).Error == nil {
					userLogged := models.UserLogged{Email: user.Email, Token: user.Token}
					c.JSON(http.StatusOK, userLogged)
					return true
				} else {
					c.JSON(http.StatusBadRequest, "User can't be saved")
				}
			} else {
				c.JSON(http.StatusBadRequest, errMap)
			}
		}
		return false
	})
}

//UserLogin serves the route POST /users/login
func UserLogin(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var userLogin models.UserLogin
		if err := c.BindJSON(&userLogin); err == nil {
			var user models.User
			if err := txn.Find(&user, "email = ?", userLogin.Email).Error; err == nil && user.Email != "" && user.LoggedIn(userLogin) {
				userLogged := models.UserLogged{Email: user.Email, Token: user.Token}
				c.JSON(http.StatusOK, userLogged)
				return true
			} else {
				c.JSON(http.StatusBadRequest, "User not found")
			}
		}
		return false
	})
}
