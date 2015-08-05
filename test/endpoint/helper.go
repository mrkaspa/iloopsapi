package endpoint

import (
	"bitbucket.org/kiloops/api/models"
	"github.com/jinzhu/gorm"
)

func saveUser() models.User {
	user := models.User{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
	models.Gdb.InTx(func(txn *gorm.DB) {
		txn.Save(&user)
	})
	return user
}
