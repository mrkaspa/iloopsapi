package endpoint

import (
	"net/http"

	"bitbucket.org/kiloops/api/ierrors"
	"bitbucket.org/kiloops/api/models"
	"gopkg.in/bluesuncorp/validator.v6"

	"github.com/gin-gonic/gin"
)

func userSession(c *gin.Context) *models.User {
	userParam, _ := c.Get("userSession")
	user := userParam.(*models.User)
	return user
}

func currentProject(c *gin.Context) *models.Project {
	projectParam, _ := c.Get("currentProject")
	project := projectParam.(*models.Project)
	return project
}

func errorResponse(c *gin.Context, errs ...error) {
	err := errs[0]
	switch err := err.(type) {
	case ierrors.AppError:
		c.JSON(http.StatusConflict, err)
	case error:
		if len(errs) == 2 {
			errorResponse(c, errs[1])
		} else {
			c.JSON(http.StatusConflict, ierrors.AppError{Code: ierrors.ErrCodeGeneral, ErrorS: err.Error()})
		}
	}
}

func errorResponseMap(c *gin.Context, errMap validator.ValidationErrors) {
	err := ierrors.AppError{Code: ierrors.ErrCodeValidation, MapErrors: make(map[string]string)}
	for _, value := range errMap {
		err.MapErrors[value.Field] = value.Tag
	}
	c.JSON(http.StatusConflict, err)
}
