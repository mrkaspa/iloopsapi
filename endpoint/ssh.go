package endpoint

import (
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//SSHCreate serves the route POST /ssh
func SSHCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var ssh models.SSH
		if err := c.BindJSON(&ssh); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return false
		}
		if valid, errMap := models.ValidStruct(&ssh); !valid {
			errorResponseMap(c, errMap)
			return false
		}
		user := userSession(c)
		ssh.UserID = user.ID
		if txn.Create(&ssh).Error != nil {
			errorResponseFromAppError(c, ErrSSHCreate)
			return false
		}
		c.JSON(http.StatusOK, ssh)
		return true
	})
}

//SSHDestroy serves the route DELETE /ssh/:id
func SSHDestroy(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var ssh models.SSH
		id, _ := strconv.Atoi(c.Param("id"))
		if txn.First(&ssh, id); ssh.ID == 0 {
			c.JSON(http.StatusNotFound, "SSH not found")
			return false
		}
		if txn.Delete(&ssh).Error != nil {
			errorResponseFromAppError(c, ErrSSHDelete)
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}
