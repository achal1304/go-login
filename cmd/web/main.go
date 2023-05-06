package main

import (
	"database/sql"
	"fmt"

	"github.com/achal1304/go-login/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	users *mysql.UserModel
}

func main() {
	dsn := "root:achal1234@/login?parseTime=true"
	db := connect(dsn)
	defer db.Close()

	app := &application{
		users: &mysql.UserModel{DB: db},
	}

	app.RunServer()
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
