package endpoint

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//Proxy requests to another api
func Proxy(c *gin.Context) {
	fmt.Println("request >> ")
	fmt.Println(c.Request.URL)
	c.Request.URL.Host = os.Getenv("GUARTZ_HOST")
	client := http.Client{}
	resp, err := client.Do(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.JSON(resp.StatusCode, string(jsonDataFromHTTP))
}
