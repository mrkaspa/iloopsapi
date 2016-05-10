package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mrkaspa/iloopsapi/endpoint"
	"github.com/mrkaspa/iloopsapi/gitadmin"
	"github.com/mrkaspa/iloopsapi/models"
	"github.com/mrkaspa/iloopsapi/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	utils.InitLog()
	utils.InitEmail()
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
