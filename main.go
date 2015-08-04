package main

import (
	"log"
	"net/http"
	"time"

	"bitbucket.org/kiloops/api/endpoint"
	"bitbucket.org/kiloops/api/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	initEnv()

	router := GetMainEngine()

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

//GetMainEngine server & routes
func GetMainEngine() *gin.Engine {
	models.InitDB()
	router := gin.Default()

	v1 := router.Group("v1")
	{
		v1.POST("users", endpoint.UserCreate)
		v1.POST("users/login", endpoint.UserLogin)
	}

	return router
}
