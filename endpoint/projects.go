package endpoint

import (
	"fmt"
	"net/http"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ProjectList serves the route GET /projects
func ProjectList(c *gin.Context) {
	user := userSession(c)
	projects := user.AllProjects()
	c.JSON(http.StatusOK, *projects)
}

//ProjectShow serves the route GET /projects/:slug
func ProjectShow(c *gin.Context) {
	c.JSON(http.StatusOK, currentProject(c))
}

//ProjectCreate serves the route POST /projects
func ProjectCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var project models.Project
		if err := c.BindJSON(&project); err != nil {
			c.JSON(http.StatusConflict, err)
			return false
		}
		if valid, errMap := models.ValidStruct(&project); !valid {
			c.JSON(http.StatusConflict, errMap)
			return false
		}
		user := userSession(c)
		if err := user.CreateProject(txn, &project); err != nil {
			c.JSON(http.StatusBadRequest, "Couldn't create the project")
			return false
		}
		c.JSON(http.StatusOK, project)
		return true
	})
}

//ProjectDestroy serves the route DELETE /projects/:slug
func ProjectDestroy(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		if err := txn.Delete(&project).Error; err != nil {
			c.JSON(http.StatusBadRequest, "Could not delete the project")
			return false
		}
		c.JSON(http.StatusOK, project)
		return true
	})
}

//ProjectLeave serves the route PUT /projects/:slug/leave
func ProjectLeave(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		user := userSession(c)
		project := currentProject(c)
		if user.HasAdminAccessTo(project.ID) {
			c.JSON(http.StatusBadRequest, "An admin user can't leave a project")
			return false
		}
		if err := user.LeaveProject(txn, project.ID); err != nil {
			c.JSON(http.StatusBadRequest, "Could not leave the project")
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}

//ProjectAddUser serves the route PUT /projects/:slug/add/:email
func ProjectAddUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		user, err := models.FindUserByEmail(c.Param("email"))
		if err != nil {
			c.JSON(http.StatusNotFound, "")
			return false
		}
		if err := project.AddUser(txn, user); err != nil {
			c.JSON(http.StatusBadRequest, "")
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}

//ProjectRemoveUser serves the route DELETE /projects/:slug/remove/:email
func ProjectRemoveUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		user, err := models.FindUserByEmail(c.Param("email"))
		if err != nil {
			c.JSON(http.StatusNotFound, "")
			return false
		}
		if err := project.RemoveUser(txn, user); err != nil {
			c.JSON(http.StatusBadRequest, "")
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}

//ProjectDelegate serves the route PUT /projects/:slug/delegate/:email
func ProjectDelegate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		userAdmin := userSession(c)
		project := currentProject(c)
		user, err := models.FindUserByEmail(c.Param("email"))
		if err != nil {
			c.JSON(http.StatusNotFound, "")
			return false
		}
		if err := project.DelegateUser(txn, userAdmin, user); err != nil {
			c.JSON(http.StatusBadRequest, "")
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}

//ProjectSchedule serves the route POST /projects/:slug/schedule
func ProjectSchedule(c *gin.Context) {
	fmt.Println("Entro en ProjectSchedule")
	models.InTx(func(txn *gorm.DB) bool {
		project, err := models.FindProjectBySlug(c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, "")
			return false
		}
		var projectConfig models.ProjectConfig
		if err := c.BindJSON(&projectConfig); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return false
		}
		project.Periodicity = projectConfig.Loops.CronFormat
		project.Command = project.GetCommand()
		if err := txn.Save(&project).Error; err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return false
		}
		if err := project.Schedule(); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}
