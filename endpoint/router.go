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

		auth := v1.Group("", Authorized())
		{
			auth.POST("ssh", SSHCreate)

			active := auth.Group("", Active())
			{
				active.DELETE("ssh/:id", SSHDestroy)
				active.GET("projects", ProjectList)
				active.GET("projects/:slug", WriteAccessToProject(), ProjectShow)
				active.POST("projects", ProjectCreate)
				active.PUT("projects/:slug/leave", WriteAccessToProject(), ProjectLeave)
				active.PUT("projects/:slug/add/:email", AdminAccessToProject(), ProjectAddUser)
				active.DELETE("projects/:slug/remove/:email", AdminAccessToProject(), ProjectRemoveUser)
				active.PUT("projects/:slug/delegate/:email", AdminAccessToProject(), ProjectDelegate)
				active.DELETE("projects/:slug", AdminAccessToProject(), ProjectDestroy)
			}
		}

		internal := v1.Group("")
		{
			internal.POST("projects/:slug/schedule", ProjectSchedule)
		}

	}

	router.NoRoute(Proxy)

	return router
}
