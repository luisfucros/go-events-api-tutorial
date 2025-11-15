package store

import (
	"context"
	"database/sql"
)

type UserStore struct {
	DB *sql.DB
}

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"username"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

func (m *UserStore) Insert(ctx context.Context, user *User) error {
	query := "INSERT INTO users (email, password, name) VALUES (?, ?, ?)"

	result, err := m.DB.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.Id = id
	return nil
}

func (m *UserStore) getUser(ctx context.Context, query string, args ...interface{}) (*User, error) {
	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserStore) Get(ctx context.Context, id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	return m.getUser(ctx, query, id)
}

func (m *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	return m.getUser(ctx, query, email)
}
