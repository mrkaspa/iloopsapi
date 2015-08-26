package endpoint

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//Proxy requests to another api
func Proxy(c *gin.Context) {
	//TODO set guartz api host
	c.Request.URL.Host = os.Getenv("GUARTZ_HOST")
	client := http.Client{}
	if resp, err := client.Do(c.Request); err == nil {
		if jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body); err == nil {
			c.JSON(resp.StatusCode, string(jsonDataFromHTTP))
			return
		}
	}
	c.JSON(http.StatusInternalServerError, "")
}
