package test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/joho/godotenv"

	"bitbucket.org/kiloops/api"
	"bitbucket.org/kiloops/api/models"
)

var apiVersion = "v1"

func around(f func(c client)) {
	ts := beforeEach()
	c := client{
		baseURL:     ts.URL + "/" + apiVersion,
		contentType: "application/json",
	}
	f(c)
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
	if err := godotenv.Load(".env_test"); err != nil {
		log.Fatal("Error loading .env_test file")
	}
}

//Client for http requests
type client struct {
	baseURL     string
	contentType string
}

func (c client) callRequest(method string, path string, reader io.Reader) (*http.Response, error) {
	return c.callRequestWithHeaders(method, path, reader, make(map[string]string))
}

func (c client) callRequestWithHeaders(method string, path string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", c.contentType)
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	return client.Do(req)
}
