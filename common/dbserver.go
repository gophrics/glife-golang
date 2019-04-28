package dbserver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Instance : Singleton Instance
var Instance *sql.DB

func openDB() *sql.DB {
	url := "mysql"
	username := "root"
	password := ""
	dbname := "glife"
	Instance, _ = sql.Open(url, username+":"+password+"@/"+dbname)
	fmt.Printf("DB OPened\n")
	return Instance
}

func init() {
	openDB()
}
