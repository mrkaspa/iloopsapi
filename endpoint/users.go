package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/ierrors"
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
		userLogged := models.UserLogged{ID: user.ID, Email: user.Email, Token: user.Token}
		c.JSON(http.StatusOK, userLogged)
	}
}
