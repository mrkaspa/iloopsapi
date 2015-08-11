package endpoint

import (
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
	var project models.Project
	models.Gdb.First(&project, id)
	if project.ID != 0 {
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

//ProjectLeave serves the route PUT /projects/:id
func ProjectLeave(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		user := userSession(c)
		id, _ := strconv.Atoi(c.Param("id"))
		if !user.HasAdminAccessTo(id) {
			if err := user.LeaveProject(txn, id); err == nil {
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "Could not leave the project")
			}
		} else {
			c.JSON(http.StatusForbidden, "")
		}
	})
}

//ProjectAddUser serves the route PUT /projects/:id
func ProjectAddUser(c *gin.Context) {

}

//ProjectDelegate serves the route PUT /projects/:id
func ProjectDelegate(c *gin.Context) {

}
