package endpoint

import (
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/validator.v2"
)

//SSHCreate serves the route POST /ssh
func SSHCreate(c *gin.Context) {
	models.Gdb.InTx(func(txn *gorm.DB) {
		var ssh models.SSH
		if err := c.BindJSON(&ssh); err == nil {
			if err := validator.Validate(&ssh); err == nil {
				user := userSession(c)
				ssh.UserID = user.ID
				if txn.Save(&ssh).Error == nil {
					c.JSON(http.StatusOK, "")
				} else {
					errorJSON(c, "SSH can't be saved")
				}
			} else {
				c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
			}
		}
	})
}

//SSHDestroy serves the route DELETE /ssh/:id
func SSHDestroy(c *gin.Context) {
	models.Gdb.InTx(func(txn *gorm.DB) {
		var ssh models.SSH
		id, _ := strconv.Atoi(c.Param("id"))
		if txn.First(&ssh, id); ssh.ID != 0 {
			if txn.Delete(&ssh).Error == nil {
				c.JSON(http.StatusOK, "")
			} else {
				errorJSON(c, "SSH can't be deleted")
			}
		} else {
			errorJSON(c, "SSH not found")
		}
	})
}
