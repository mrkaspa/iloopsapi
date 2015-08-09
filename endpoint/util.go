package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/models"

	"github.com/gin-gonic/gin"
)

func errorJSON(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, map[string]string{"error": err})
}

func userSession(c *gin.Context) *models.User {
	userParam, _ := c.Get("userSession")
	user := userParam.(models.User)
	return &user
}
