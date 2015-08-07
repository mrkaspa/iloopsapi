package endpoint

import "github.com/gin-gonic/gin"

func ProjectList(c *gin.Context) {

}

func ProjectShow(c *gin.Context) {

}

func ProjectCreate(c *gin.Context) {
	// models.Gdb.InTx(func(txn *gorm.DB) {
	// 	var project models.Project
	// 	if err := c.BindJSON(&project); err == nil {
	// 		if err := validator.Validate(&project); err == nil {
	// 			if txn.Save(&project).Error == nil {
	//         userProject := models.UsersProjects{Role: models.Creator}
	// 				c.JSON(http.StatusOK, "")
	// 			} else {
	// 				errorJSON(c, "Project can't be saved")
	// 			}
	// 		} else {
	// 			c.JSON(http.StatusBadRequest, err.(validator.ErrorMap))
	// 		}
	// 	}
	// })
}

func ProjectDestroy(c *gin.Context) {

}
