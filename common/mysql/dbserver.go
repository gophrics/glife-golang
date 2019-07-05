package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Importing mysql connector for golang
)

// Instance : Singleton Instance
var Instance *sql.DB

func openDB() *sql.DB {
	url := "mysql"
	username := "root"
	password := ""
	dbname := "glife"
	Instance, _ = sql.Open(url, username+":"+password+"@/"+dbname)
	return Instance
}

func init() {
	openDB()
}
