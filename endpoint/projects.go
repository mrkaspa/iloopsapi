package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/utils"
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
			c.JSON(http.StatusBadRequest, err)
			return false
		}
		if valid, errMap := models.ValidStruct(&project); !valid {
			errorResponseMap(c, errMap)
			return false
		}
		user := userSession(c)
		if err := user.CreateProject(txn, &project); err != nil {
			if err == models.ErrUserExceedMaxProjects {
				errorResponse(c, http.StatusConflict, err)
			} else {
				errorResponseFromAppError(c, ErrProjectCreate)
			}
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
			errorResponseFromAppError(c, ErrProjectDelete)
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
			errorResponseFromAppError(c, ErrAdminCantLeaveProject)
			return false
		}
		if err := user.LeaveProject(txn, project.ID); err != nil {
			errorResponseFromAppError(c, ErrUserLeaveProject)
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
			if err == models.ErrUserExceedMaxProjects {
				errorResponse(c, http.StatusConflict, err)
			} else {
				errorResponseFromAppError(c, ErrProjectAddUser)
			}
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
			errorResponseFromAppError(c, ErrProjectRemoveUser)
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
			errorResponseFromAppError(c, ErrProjectDelegateUser)
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}

//ProjectSchedule serves the route POST /projects/:slug/schedule
func ProjectSchedule(c *gin.Context) {
	utils.Log.Info("Entro en ProjectSchedule")
	models.InTx(func(txn *gorm.DB) bool {
		project, err := models.FindProjectBySlug(c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, "")
			return false
		}
		var projectConfig models.ProjectConfig
		if err := c.BindJSON(&projectConfig); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return false
		}
		project.Periodicity = projectConfig.Loops.CronFormat
		project.Command = project.GetCommand()
		if err := txn.Save(&project).Error; err != nil {
			errorResponse(c, ErrCodeGeneral, err)
			return false
		}
		if err := project.Schedule(); err != nil {
			errorResponse(c, ErrCodeGeneral, err)
			return false
		}
		c.JSON(http.StatusOK, "")
		return true
	})
}
