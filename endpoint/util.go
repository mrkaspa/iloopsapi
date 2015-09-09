package endpoint

import (
	"net/http"

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

func errorResponseFromAppError(c *gin.Context, error AppError) {
	errorResponse(c, error.Code, error.Error)
}

func errorResponse(c *gin.Context, code int, err error) {
	c.JSON(http.StatusConflict, JSONError{Code: code, ErrorCad: err.Error()})
}

func errorResponseMap(c *gin.Context, errMap validator.ValidationErrors) {
	jsonErr := JSONError{Code: ErrCodeValidation, MapErrors: make(map[string]string)}
	for _, value := range errMap {
		jsonErr.MapErrors[value.Field] = value.Tag
	}
	c.JSON(http.StatusConflict, jsonErr)
}
