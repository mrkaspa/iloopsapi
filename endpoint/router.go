package endpoint

import "github.com/gin-gonic/gin"

//GetMainEngine server & routes
func GetMainEngine() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	{
		v1.POST("users", UserCreate)
		v1.POST("users/login", UserLogin)
		v1.POST("ssh", SSHCreate)
		v1.DELETE("ssh/:id", SSHDestroy)
	}

	return router
}
