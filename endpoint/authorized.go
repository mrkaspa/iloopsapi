package endpoint

import (
	"net/http"

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
				c.Set("userSession", &user)
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

//AdminAccessToProject middleware
func AdminAccessToProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := userSession(c)
		if project, err := models.FindProjectBySlug(c.Param("slug")); err == nil {
			if user.HasAdminAccessTo(project.ID) {
				c.Set("currentProject", project)
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

//WriteAccessToProject middleware
func WriteAccessToProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := userSession(c)
		if project, err := models.FindProjectBySlug(c.Param("slug")); err == nil {
			if user.HasWriteAccessTo(project.ID) {
				c.Set("currentProject", project)
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}
