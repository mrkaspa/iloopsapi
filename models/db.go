package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

//KDB database connection
type KDB struct {
	*gorm.DB
}

//Gdb connection
var Gdb KDB

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

	//Add unique index
	db.Model(&User{}).AddUniqueIndex("idx_user_email", "email")
	db.Model(&SSH{}).AddUniqueIndex("idx_ssh_hash", "hash")
	db.Model(&Project{}).AddUniqueIndex("idx_project_slug", "slug")
	db.Model(&UsersProjects{}).AddUniqueIndex("idx_user_project", "user_id", "project_id")

	//Add FK
	// db.Model(&SSH{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	// db.Model(&UsersProjects{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	// db.Model(&UsersProjects{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")

	Gdb = KDB{&db}
}

//InTx executes function in a transaction
func (kdb KDB) InTx(f func(*gorm.DB)) {
	txn := kdb.InitTx()
	defer Commit(txn)
	f(txn)
}

//InitTx creates a transaction
func (kdb KDB) InitTx() *gorm.DB {
	txn := kdb.Begin()
	if txn.Error != nil {
		fmt.Println(txn.Error)
		panic(txn.Error)
	}
	return txn
}

//Commit a transaction
func Commit(txn *gorm.DB) {
	txn.Commit()
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

//Rollback a transaction
func Rollback(txn *gorm.DB) {
	txn.Rollback()
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}
