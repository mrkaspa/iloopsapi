package endpoint

import (
	"bitbucket.org/kiloops/api/models"
	"github.com/jinzhu/gorm"
)

func inTxn(f func(txn *gorm.DB)) {
	txn := models.InitTx()
	defer models.Commit(txn)
	f(txn)
}
