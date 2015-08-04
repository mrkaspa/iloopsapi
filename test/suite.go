package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/joho/godotenv"

	"bitbucket.org/kiloops/api"
	"bitbucket.org/kiloops/api/models"
)

var apiVersion = "v1"

func Around(f func(c *Client)) {
	ts := beforeEach()
	c := Client{
		&http.Client{},
		ts.URL + "/" + apiVersion,
		"application/json",
	}
	f(&c)
	afterEach(ts)
}

func beforeEach() *httptest.Server {
	initEnv()
	models.InitDB()
	ts := httptest.NewServer(main.GetMainEngine())
	return ts
}

func afterEach(ts *httptest.Server) {
	models.Gdb.Delete(models.User{})
	ts.Close()
}

func initEnv() {
	path := ".env_test"
	for i := 1; ; i++ {
		if err := godotenv.Load(path); err != nil {
			if i > 3 {
				panic("Error loading .env_test file")
			} else {
				path = "../" + path
			}
		} else {
			break
		}
	}
}

//Client for http requests
type Client struct {
	*http.Client
	baseURL     string
	contentType string
}

func (c Client) CallRequest(method string, path string, reader io.Reader) (*http.Response, error) {
	return c.CallRequestWithHeaders(method, path, reader, make(map[string]string))
}

func (c Client) CallRequestWithHeaders(method string, path string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, c.baseURL+path, reader)
	req.Header.Set("Content-Type", c.contentType)
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	return c.Do(req)
}
