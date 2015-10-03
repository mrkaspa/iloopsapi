package endpoint

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

//Proxy requests to another api
func Proxy(c *gin.Context) {
	c.Request.RequestURI = ""
	newURLString := fmt.Sprintf("http://%s%s?%s", os.Getenv("GUARTZ_HOST"), c.Request.URL.Path, c.Request.URL.RawQuery)
	newURL, err := url.Parse(newURLString)

	c.Request.URL = newURL
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
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	c.Writer.Write(jsonDataFromHTTP)
}
