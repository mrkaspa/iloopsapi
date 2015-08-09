package endpoint

import "github.com/gin-gonic/gin"

//GetMainEngine server & routes
func GetMainEngine() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	{
		v1.POST("users", UserCreate)
		v1.POST("users/login", UserLogin)

		v1.POST("ssh", Authorized(), SSHCreate)
		v1.DELETE("ssh/:id", Authorized(), SSHDestroy)

		v1.GET("projects", Authorized(), ProjectList)
		v1.GET("projects/:id", Authorized(), WriteAccessToProject(), ProjectShow)
		v1.POST("projects", Authorized(), ProjectCreate)
		v1.DELETE("projects/:id", Authorized(), AdminAccessToProject(), ProjectDestroy)
	}

	return router
}
