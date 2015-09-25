package endpoint

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StatusGet(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]int64{"timestamp": time.Now().Unix()})
}
