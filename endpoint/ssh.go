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
		if err := c.BindJSON(&ssh); err == nil {
			if valid, errMap := models.ValidStruct(&ssh); valid {
				user := userSession(c)
				ssh.UserID = user.ID
				if txn.Create(&ssh).Error == nil {
					c.JSON(http.StatusOK, ssh)
					return true
				} else {
					c.JSON(http.StatusBadRequest, "SSH can't be saved")
				}
			} else {
				c.JSON(http.StatusConflict, errMap)
			}
		}
		return false
	})
}

//SSHDestroy serves the route DELETE /ssh/:id
func SSHDestroy(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var ssh models.SSH
		id, _ := strconv.Atoi(c.Param("id"))
		if txn.First(&ssh, id); ssh.ID != 0 {
			if txn.Delete(&ssh).Error == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				c.JSON(http.StatusBadRequest, "SSH can't be deleted")
			}
		} else {
			c.JSON(http.StatusNotFound, "SSH not found")
		}
		return false
	})
}
