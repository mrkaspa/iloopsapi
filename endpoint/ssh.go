package endpoint

import (
	"net/http"
	"strconv"

	"github.com/mrkaspa/iloopsapi/ierrors"
	"github.com/mrkaspa/iloopsapi/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

//SSHList serves the route GET /ssh
func SSHList(c *gin.Context) {
	user := userSession(c)
	sshs := user.AllSHHs()
	c.JSON(http.StatusOK, *sshs)
}

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
		ssh.Name = slug.Make(ssh.Name)
		ssh.UserID = user.ID
		if err := txn.Create(&ssh).Error; err != nil {
			errorResponse(c, err, ierrors.ErrSSHCreate)
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
		strID := c.Param("id")
		id, _ := strconv.Atoi(strID)
		user := userSession(c)
		txn.Where("(id = ? or name like ?) and user_id = ?", id, strID, user.ID).First(&ssh)
		if ssh.ID == 0 {
			c.JSON(http.StatusNotFound, "SSH not found")
			return false
		}
		if err := txn.Delete(&ssh).Error; err != nil {
			errorResponse(c, err, ierrors.ErrSSHDelete)
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}
