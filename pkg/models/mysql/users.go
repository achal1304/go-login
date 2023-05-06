package mysql

import (
	"database/sql"
	"fmt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) insert(name string, email string, address string, password string) error {
	fmt.Print("Inside insert")
	return nil
}
