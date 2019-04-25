package dbserver

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func openDB() *sql.DB {
	url := "mysql"
	username := "root"
	password := ""
	dbname := "glife"
	db, err := sql.Open(url, username+":"+password+"@/"+dbname)

	if err != nil {
		panic(err.Error())
	}
	return db
}

func close() {
	Instance.Close()
}

func init() {
	openDB()
}

// Instance : Singleton Instance
var Instance = openDB()
