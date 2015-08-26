package main

import (
	"log"
	"net/http"
	"time"

	"bitbucket.org/kiloops/api/endpoint"
	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	models.InitDB()
	gitadmin.InitVars()
	gitadmin.InitGitAdmin()

	router := endpoint.GetMainEngine()

	s := &http.Server{
		Addr:         ":8080",
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
