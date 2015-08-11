package endpoint

import (
	"bitbucket.org/kiloops/api/models"

	"github.com/gin-gonic/gin"
)

func userSession(c *gin.Context) *models.User {
	userParam, _ := c.Get("userSession")
	user := userParam.(models.User)
	return &user
}
