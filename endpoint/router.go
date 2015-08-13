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
			auth.POST("projects", ProjectCreate)
			auth.GET("projects/:id", WriteAccessToProject(), ProjectShow)
			auth.PUT("projects/:id/leave", WriteAccessToProject(), ProjectLeave)
			auth.PUT("projects/:id/add/:user_id", AdminAccessToProject(), ProjectAddUser)
			auth.PUT("projects/:id/delegate/:user_id", AdminAccessToProject(), ProjectDelegate)
			auth.DELETE("projects/:id", AdminAccessToProject(), ProjectDestroy)
		}

		internal := v1.Group("")
		{
			internal.POST("/executions/:project_id", ExecutionCreate)
		}
	}

	return router
}
