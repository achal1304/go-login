package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail     = errors.New("models: email address already in use")
	ErrInvalidCredentials = errors.New("models: invalid user credentials")
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email string, password string) error {
	fmt.Print("Inside insert")

	hashedpassowrd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO USER(email, hashed_password, createdAt) VALUES(?,?,UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, email, hashedpassowrd)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if strings.Contains(mysqlError.Message, "user_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) AuthenticateUser(email string, password string) (int, error) {
	var id int
	var hashed_password_user []byte
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	stmt := `SELECT userId, hashed_password from user where email = ?`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashed_password_user)
	fmt.Print("hashed password - ", hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
