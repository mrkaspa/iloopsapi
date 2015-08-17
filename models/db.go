package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

//Gdb connection
var Gdb *gorm.DB

//InitDB connection
func InitDB() {
	//open db
	fmt.Println("*** INIT DB ***")
	//connString := revel.Config.StringDefault("db.conn", "")
	connString := os.Getenv("MYSQL_DB")
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		fmt.Println("Unable to connect to the database")
		revel.ERROR.Println("FATAL", err)
		panic(err)
	}
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//Migrations
	db.AutoMigrate(&User{})
	db.AutoMigrate(&SSH{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&UsersProjects{})
	db.AutoMigrate(&Execution{})

	//Add unique index
	db.Model(&UsersProjects{}).AddUniqueIndex("idx_user_project", "user_id", "project_id")

	//Add FK
	db.Model(&SSH{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&UsersProjects{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&UsersProjects{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	db.Model(&Execution{}).AddForeignKey("project_id", "projects(id)", "RESTRICT", "RESTRICT")

	Gdb = &db
}

//InTx executes function in a transaction
func InTx(f func(*gorm.DB) bool) {
	txn := Gdb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	if f(txn) == true {
		txn.Commit()
	} else {
		txn.Rollback()
	}
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}
