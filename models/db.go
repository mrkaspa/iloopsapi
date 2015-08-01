package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// it can be used for jobs
var Gdb gorm.DB

// init db
func InitDB() {
	// open db
	fmt.Println("*** INIT DB ***")
	// connString := revel.Config.StringDefault("db.conn", "")
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
	db.AutoMigrate(&User{})

	// Add unique index
	db.Model(&User{}).AddUniqueIndex("idx_user_email", "email")
	Gdb = db
}

func InitTx() *gorm.DB {
	txn := Gdb.Begin()
	if txn.Error != nil {
		fmt.Println("Unable to open a connection")
		fmt.Println(txn.Error)
		panic(txn.Error)
	}
	return txn
}

// This method clears the c.Txn after each transaction
func Commit(txn *gorm.DB) {
	txn.Commit()
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

// This method clears the c.Txn after each transaction, too
func Rollback(txn *gorm.DB) {
	txn.Rollback()
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}
