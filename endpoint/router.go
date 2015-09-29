package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//GetMainEngine server & routes
func GetMainEngine() *gin.Engine {
	router := gin.Default()

	root := router.Group("", AngularFilter)

	root.OPTIONS("/*path", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	v1 := root.Group("v1")
	{
		v1.GET("status", StatusGet)

		v1.POST("users", UserCreate)
		v1.POST("users/login", UserLogin)

		v1.POST("ssh", SSHCreate)
		v1.DELETE("ssh/:id", SSHDestroy)

		auth := v1.Group("", Authorized())
		{
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
			internal.POST("projects/:slug/schedule", ProjectSchedule)
		}

	}

	router.NoRoute(Proxy)

	return router
}
