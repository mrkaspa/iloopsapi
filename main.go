package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/kiloops/api/endpoint"
	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	utils.InitLog()
	models.InitDB()
	models.InitGuartzClient()
	gitadmin.InitVars()
	gitadmin.InitGitAdmin()

	router := endpoint.GetMainEngine()
	router.Use(gin.LoggerWithWriter(utils.LogWriter))
	router.Use(gin.Recovery())
	port := os.Getenv("PORT")

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.ListenAndServe()
}

func initEnv() {
	if err := godotenv.Load(".env_dev"); err != nil {
		log.Fatal("Error loading .env_dev file")
	}
}
