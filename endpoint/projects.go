package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/validator.v2"
)

//ProjectList serves the route GET /projects
func ProjectList(c *gin.Context) {
	user := userSession(c)
	projects := user.AllProjects()
	c.JSON(http.StatusOK, projects)
}

//ProjectShow serves the route GET /projects/:id
func ProjectShow(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if project, err := models.FindProject(id); err == nil {
		c.JSON(http.StatusOK, project)
	} else {
		c.JSON(http.StatusNotFound, "")
	}
}

//ProjectCreate serves the route POST /projects
func ProjectCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var project models.Project
		if err := c.BindJSON(&project); err == nil {
			if err := validator.Validate(&project); err == nil {
				user := userSession(c)
				if err := user.CreateProject(txn, &project); err == nil {
					c.JSON(http.StatusOK, project)
					return true
				} else {
					c.JSON(http.StatusBadRequest, "Couldn't create the project")
				}
			} else {
				c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
			}
		}
		return false
	})
}

//ProjectDestroy serves the route DELETE /projects/:id
func ProjectDestroy(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		id, _ := strconv.Atoi(c.Param("id"))
		var project models.Project
		txn.First(&project, id)
		if project.ID != 0 {
			if err := txn.Delete(&project).Error; err == nil {
				// project.DeleteRels(txn)
				c.JSON(http.StatusOK, project)
				return true
			} else {
				c.JSON(http.StatusBadRequest, "Could not delete the project")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ProjectLeave serves the route PUT /projects/:id/leave
func ProjectLeave(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		user := userSession(c)
		id, _ := strconv.Atoi(c.Param("id"))
		if !user.HasAdminAccessTo(id) {
			if err := user.LeaveProject(txn, id); err == nil {
				c.JSON(http.StatusOK, "")
				return true
			} else {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, "Could not leave the project")
			}
		} else {
			c.JSON(http.StatusForbidden, "")
		}
		return false
	})
}

//ProjectAddUser serves the route PUT /projects/:id/add/:user_id
func ProjectAddUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		id, _ := strconv.Atoi(c.Param("id"))
		userID, _ := strconv.Atoi(c.Param("user_id"))
		if project, err := models.FindProject(id); err == nil {
			if user, err := models.FindUser(userID); err == nil {
				if err := project.AddUser(txn, user); err == nil {
					c.JSON(http.StatusOK, "")
					return true
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusNotFound, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ProjectRemoveUser serves the route DELETE /projects/:id/remove/:user_id
func ProjectRemoveUser(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		id, _ := strconv.Atoi(c.Param("id"))
		userID, _ := strconv.Atoi(c.Param("user_id"))
		if project, err := models.FindProject(id); err == nil {
			if user, err := models.FindUser(userID); err == nil {
				if err := project.RemoveUser(txn, user); err == nil {
					c.JSON(http.StatusOK, "")
					return true
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusNotFound, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ProjectDelegate serves the route PUT /projects/:id/delegate/:user_id
func ProjectDelegate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		id, _ := strconv.Atoi(c.Param("id"))
		userAdmin := userSession(c)
		userID, _ := strconv.Atoi(c.Param("user_id"))
		if project, err := models.FindProject(id); err == nil {
			if user, err := models.FindUser(userID); err == nil {
				if err := project.DelegateUser(txn, userAdmin, user); err == nil {
					c.JSON(http.StatusOK, "")
					return true
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusNotFound, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//SSHHasAccess serves the route GET /projects/:id/has_access
func ProjectHasAccessBySSH(c *gin.Context) {
	var ssh models.SSH
	if err := c.BindJSON(&ssh); err == nil {
		ssh.Hash = utils.MD5(ssh.PublicKey)
		models.Gdb.Where("hash like ?", ssh.Hash).First(&ssh)
		id, _ := strconv.Atoi(c.Param("id"))
		project, err := models.FindProject(id)
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
