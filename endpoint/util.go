package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorJSON(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, map[string]string{"error": err})
}
