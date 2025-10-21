package database

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id        int64    `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"username"`
	Password  string   `json:"-"`
	CreatedAt string   `json:"created_at"`
}