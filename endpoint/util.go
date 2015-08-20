package endpoint

import (
	"bitbucket.org/kiloops/api/models"

	"github.com/gin-gonic/gin"
)

func userSession(c *gin.Context) *models.User {
	userParam, _ := c.Get("userSession")
	user := userParam.(*models.User)
	return user
}

func currentProject(c *gin.Context) *models.Project {
	projectParam, _ := c.Get("currentProject")
	project := projectParam.(*models.Project)
	return project
}
