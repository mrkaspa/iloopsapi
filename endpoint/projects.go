package endpoint

import (
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/models"
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
	models.Gdb.InTx(func(txn *gorm.DB) {
		var project models.Project
		if err := c.BindJSON(&project); err == nil {
			if err := validator.Validate(&project); err == nil {
				user := userSession(c)
				if err := user.CreateProject(txn, &project); err == nil {
					c.JSON(http.StatusOK, project)
				} else {
					errorJSON(c, err.Error())
				}
			} else {
				c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
			}
		}
	})
}

//ProjectDestroy serves the route DELETE /projects/:id
func ProjectDestroy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var project models.Project
	models.Gdb.First(&project, id)
	if project.ID != 0 {
		if models.Gdb.Delete(&project).Error != nil {
			c.JSON(http.StatusOK, project)
		} else {
			c.JSON(http.StatusBadRequest, "Could not delete the project")
		}
	} else {
		c.JSON(http.StatusNotFound, "")
	}
}

//ProjectLeave serves the route PUT /projects/:id
func ProjectLeave(c *gin.Context) {

}

//ProjectAddUser serves the route PUT /projects/:id
func ProjectAddUser(c *gin.Context) {

}

//ProjectDelegate serves the route PUT /projects/:id
func ProjectDelegate(c *gin.Context) {

}
