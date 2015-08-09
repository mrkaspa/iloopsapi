package endpoint

import (
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"

	"github.com/gin-gonic/gin"
)

//Authorized middleware
func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		authID := c.Request.Header.Get("AUTH_ID")
		authToken := c.Request.Header.Get("AUTH_TOKEN")
		if authID == "" || authToken == "" {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			var user models.User
			models.Gdb.Where("id = ? and token = ?", authID, authToken).First(&user)
			if user.ID != 0 {
				c.Set("userSession", user)
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func AdminAccessToProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := userSession(c)
		projectID, _ := strconv.Atoi(c.Param("id"))
		if user.HasAdminAccessTo(projectID) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func WriteAccessToProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := userSession(c)
		projectID, _ := strconv.Atoi(c.Param("id"))
		if user.HasWriteAccessTo(projectID) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
