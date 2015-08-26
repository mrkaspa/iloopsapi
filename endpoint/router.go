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
			auth.GET("projects/:slug", WriteAccessToProject(), ProjectShow)
			auth.PUT("projects/:slug/leave", WriteAccessToProject(), ProjectLeave)
			auth.PUT("projects/:slug/add/:email", AdminAccessToProject(), ProjectAddUser)
			auth.DELETE("projects/:slug/remove/:email", AdminAccessToProject(), ProjectRemoveUser)
			auth.PUT("projects/:slug/delegate/:email", AdminAccessToProject(), ProjectDelegate)
			auth.DELETE("projects/:slug", AdminAccessToProject(), ProjectDestroy)
		}

		internal := v1.Group("")
		{
			internal.GET("projects/:slug/has_access", ProjectHasAccessBySSH)
		}
	}

	router.NoRoute(Proxy)

	return router
}
