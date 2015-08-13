package endpoint

import (
	"net/http"
	"strconv"

	"gopkg.in/validator.v2"

	"bitbucket.org/kiloops/api/models"
	"github.com/gin-gonic/gin"
)

//ExecutionCreate serves the route POST /executions/:project_id
func ExecutionCreate(c *gin.Context) {
	models.Gdb.InTx(func(txn *models.KDB) {
		projectID, _ := strconv.Atoi(c.Param("project_id"))
		if _, err := models.FindProject(projectID); err == nil {
			var execution models.Execution
			if err := c.BindJSON(&execution); err == nil {
				if err := validator.Validate(&execution); err == nil {
					execution.ProjectID = projectID
					if txn.Save(&execution).Error == nil {
						c.JSON(http.StatusOK, "")
					} else {
						c.JSON(http.StatusBadRequest, "Execution can't be saved")
					}
				} else {
					c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
				}
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
	})
}
