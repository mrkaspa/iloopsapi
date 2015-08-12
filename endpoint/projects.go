package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
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
	models.Gdb.InTx(func(txn *models.KDB) {
		var project models.Project
		if err := c.BindJSON(&project); err == nil {
			if err := validator.Validate(&project); err == nil {
				user := userSession(c)
				if err := user.CreateProject(txn, &project); err == nil {
					c.JSON(http.StatusOK, project)
				} else {
					c.JSON(http.StatusBadRequest, "Couldn't create the project")
				}
			} else {
				c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
			}
		}
	})
}

//ProjectDestroy serves the route DELETE /projects/:id
func ProjectDestroy(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		id, _ := strconv.Atoi(c.Param("id"))
		var project models.Project
		txn.First(&project, id)
		if project.ID != 0 {
			if err := txn.Delete(&project).Error; err == nil {
				project.DeleteRels(txn)
				c.JSON(http.StatusOK, project)
			} else {
				c.JSON(http.StatusBadRequest, "Could not delete the project")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
	})
}

//ProjectLeave serves the route PUT /projects/:id/leave
func ProjectLeave(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		user := userSession(c)
		id, _ := strconv.Atoi(c.Param("id"))
		if !user.HasAdminAccessTo(id) {
			if err := user.LeaveProject(txn, id); err == nil {
				c.JSON(http.StatusOK, "")
			} else {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, "Could not leave the project")
			}
		} else {
			c.JSON(http.StatusForbidden, "")
		}
	})
}

//ProjectAddUser serves the route PUT /projects/:id/add/:user_id
func ProjectAddUser(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		id, _ := strconv.Atoi(c.Param("id"))
		userID, _ := strconv.Atoi(c.Param("user_id"))
		if project, err := models.FindProject(id); err == nil {
			if user, err := models.FindUser(userID); err == nil {
				if err := project.AddUser(txn, user); err == nil {
					c.JSON(http.StatusOK, "")
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusNotFound, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
	})
}

//ProjectDelegate serves the route PUT /projects/:id/delegate/:user_id
func ProjectDelegate(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		id, _ := strconv.Atoi(c.Param("id"))
		userAdmin := userSession(c)
		userID, _ := strconv.Atoi(c.Param("user_id"))
		if project, err := models.FindProject(id); err == nil {
			if user, err := models.FindUser(userID); err == nil {
				if err := project.DelegateUser(txn, userAdmin, user); err == nil {
					c.JSON(http.StatusOK, "")
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusNotFound, "")
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
	})
}
