package endpoint

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StatusGet(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"timestamp": time.Now().Format(time.RFC850)})
}
