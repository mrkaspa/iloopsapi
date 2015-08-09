package endpoint

import "github.com/gin-gonic/gin"

//GetMainEngine server & routes
func GetMainEngine() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	{
		v1.POST("users", UserCreate)
		v1.POST("users/login", UserLogin)

		auth := v1.Group("", Authorized())
		{
			auth.POST("ssh", SSHCreate)
			auth.DELETE("ssh/:id", SSHDestroy)

			auth.GET("projects", ProjectList)
			auth.GET("projects/:id", WriteAccessToProject(), ProjectShow)
			auth.POST("projects", ProjectCreate)
			auth.DELETE("projects/:id", AdminAccessToProject(), ProjectDestroy)
		}
	}

	return router
}
