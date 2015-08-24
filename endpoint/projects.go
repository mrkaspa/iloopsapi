package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mrkaspa/go-helpers"
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
		if err := c.BindJSON(&project); err == nil {
			if valid, errMap := models.ValidStruct(&project); valid {
				user := userSession(c)
				if err := user.CreateProject(txn, &project); err == nil {
					c.JSON(http.StatusOK, project)
					return true
				} else {
					c.JSON(http.StatusBadRequest, "Couldn't create the project")
				}
			} else {
				c.JSON(http.StatusConflict, errMap)
			}
		}
		return false
	})
}

//ProjectDestroy serves the route DELETE /projects/:slug
func ProjectDestroy(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		if err := txn.Delete(&project).Error; err == nil {
			// project.DeleteRels(txn)
			c.JSON(http.StatusOK, project)
			return true
		} else {
			c.JSON(http.StatusBadRequest, "Could not delete the project")
		}
		return false
	})
}

//ProjectLeave serves the route PUT /projects/:slug/leave
func ProjectLeave(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		user := userSession(c)
		project := currentProject(c)
		if !user.HasAdminAccessTo(project.ID) {
			if err := user.LeaveProject(txn, project.ID); err == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				c.JSON(http.StatusBadRequest, "Could not leave the project")
			}
		} else {
			c.JSON(http.StatusBadRequest, "An admin user can't leave a project")
		}
		return false
	})
}

//ProjectAddUser serves the route PUT /projects/:slug/add/:email
func ProjectAddUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		if user, err := models.FindUserByEmail(c.Param("email")); err == nil {
			if err := project.AddUser(txn, user); err == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ProjectRemoveUser serves the route DELETE /projects/:slug/remove/:email
func ProjectRemoveUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		project := currentProject(c)
		if user, err := models.FindUserByEmail(c.Param("email")); err == nil {
			if err := project.RemoveUser(txn, user); err == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ProjectDelegate serves the route PUT /projects/:slug/delegate/:email
func ProjectDelegate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		userAdmin := userSession(c)
		project := currentProject(c)
		if user, err := models.FindUserByEmail(c.Param("email")); err == nil {
			if err := project.DelegateUser(txn, userAdmin, user); err == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//SSHHasAccess serves the route GET /projects/:slug/has_access
func ProjectHasAccessBySSH(c *gin.Context) {
	var ssh models.SSH
	if err := c.BindJSON(&ssh); err == nil {
		ssh.Hash = helpers.MD5(ssh.PublicKey)
		models.Gdb.Where("hash like ?", ssh.Hash).First(&ssh)
		project, err := models.FindProjectBySlug(c.Param("slug"))
		if ssh.ID != 0 && err == nil {
			var user models.User
			models.Gdb.Model(&ssh).Related(&user)
			if user.HasWriteAccessTo(project.ID) {
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusForbidden, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
	}
}
