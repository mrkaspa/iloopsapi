package models

import (
	"database/sql"
	"os"

	"bitbucket.org/kiloops/api/utils"

	"github.com/jinzhu/gorm"
)

//Gdb connection
var Gdb *gorm.DB

//InitDB connection
func InitDB() {
	//open db
	utils.Log.Info("*** INIT DB ***")
	connString := os.Getenv("MYSQL_DB")
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		utils.Log.Panic(err)
	}
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//Migrations
	db.AutoMigrate(&User{})
	db.AutoMigrate(&SSH{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&UsersProjects{})

	//Add unique index
	db.Model(&UsersProjects{}).AddUniqueIndex("idx_user_project", "user_id", "project_id")

	//Add FK
	db.Model(&SSH{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&UsersProjects{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&UsersProjects{}).AddForeignKey("project_id", "projects(id)", "RESTRICT", "RESTRICT")

	Gdb = &db
}

//InTx executes function in a transaction
func InTx(f func(*gorm.DB) bool) {
	utils.Log.Info("***INIT TRANSACTION***")
	txn := Gdb.Begin()
	if txn.Error != nil {
		utils.Log.Panic(txn.Error)
	}
	if f(txn) == true {
		utils.Log.Info("***TRANSACTION COMMITED***")
		txn.Commit()
	} else {
		utils.Log.Info("***TRANSACTION ROLLBACK***")
		txn.Rollback()
	}
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		utils.Log.Panic(err)
	}
}
