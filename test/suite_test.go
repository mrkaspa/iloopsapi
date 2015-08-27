package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"os"

	"github.com/joho/godotenv"

	"bitbucket.org/kiloops/api/endpoint"
	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/utils"

	gEndpoint "github.com/infiniteloopsco/guartz/endpoint"
	gModels "github.com/infiniteloopsco/guartz/models"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts         *httptest.Server
	gServer    *httptest.Server
	client     utils.Client
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
	gModels.InitDB()
	gitadmin.InitVars()
	gitadmin.InitGitAdmin()
	gModels.InitCron()
	cleanDB()
	ts = httptest.NewServer(endpoint.GetMainEngine())
	gServer = httptest.NewServer(gEndpoint.GetMainEngine())
	gURL, _ := url.Parse(gServer.URL)
	os.Setenv("GUARTZ_HOST", gURL.Host)
	client = utils.Client{
		&http.Client{},
		ts.URL + "/" + apiVersion,
		"application/json",
	}
})

var _ = AfterSuite(func() {
	models.Gdb.Close()
	gModels.Gdb.Close()
	ts.Close()
	gServer.Close()
	gitadmin.FinishGitAdmin()
})

var _ = BeforeEach(func() {
	cleanDB()
})

func cleanDB() {
	fmt.Println("***Cleaning***")
	models.Gdb.Delete(models.UsersProjects{})
	models.Gdb.Delete(models.Project{})
	models.Gdb.Delete(models.SSH{})
	models.Gdb.Delete(models.User{})
	gModels.Gdb.Delete(gModels.Execution{})
	gModels.Gdb.Delete(gModels.Task{})
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
