package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//UserCreate serves the route POST /users
func UserCreate(c *gin.Context) {
	inTxn(func(txn *gorm.DB) {
		var user models.User
		if err := c.BindJSON(&user); err == nil {
			if txn.Save(&user).Error == nil {
				c.JSON(http.StatusOK, user)
				return
			}
			c.JSON(http.StatusBadRequest, user)
		}
	})
}
