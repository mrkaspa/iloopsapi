package endpoint

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

	"bitbucket.org/kiloops/api/endpoint"
	"bitbucket.org/kiloops/api/models"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts         *httptest.Server
	client     Client
	apiVersion = "v1"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	fmt.Println("Suite found")
	RunSpecs(t, "Api Suite")
}

var _ = BeforeSuite(func() {
	initEnv()
	models.InitDB()
	ts = httptest.NewServer(endpoint.GetMainEngine())
	client = Client{
		&http.Client{},
		ts.URL + "/" + apiVersion,
		"application/json",
	}
})

var _ = AfterSuite(func() {
	// cleanDB()
	models.Gdb.Close()
	ts.Close()
})

func cleanDB() {
	fmt.Println("***Cleaning***")
	models.Gdb.Delete(models.User{})
	models.Gdb.Delete(models.SSH{})
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
