package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/achal1304/go-login/internal/mailer"
	"github.com/achal1304/go-login/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	users   *mysql.UserModel
	session *sessions.Session
	mailer  mailer.Mailer
}

func main() {
	dsn := "root:achal1234@/login?parseTime=true"
	secret := "s6Nd%+pPbnzHbS*+9Pk8qGWhTzbpa@ge"

	db := connect(dsn)
	defer db.Close()

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		users:   &mysql.UserModel{DB: db},
		session: session,
		mailer:  mailer.New(),
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
