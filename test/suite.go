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

func around(f func(c client)) {
	ts := beforeEach()
	c := client{BaseURL: ts.URL}
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
	BaseURL string
}

func (c client) callRequest(method string, path string, contentType string, reader io.Reader) (*http.Response, error) {
	return c.callRequestWithHeaders(method, path, contentType, reader, make(map[string]string))
}

func (c client) callRequestWithHeaders(method string, path string, contentType string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	return client.Do(req)
}
