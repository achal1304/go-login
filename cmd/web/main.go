package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:achal1234@/login?parseTime=true"
	connect(dsn)

	RunServer()
	fmt.Println("hello world")
}

func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		print(err)
	}
	return db
}
