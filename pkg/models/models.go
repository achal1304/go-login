package models

import "time"

type User struct {
	UserId          int
	Name            string
	Address         string
	Email           string
	Hashed_password []byte
	CreatedAt       time.Time
}
