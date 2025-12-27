package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserInsert_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}
	u := &User{Email: "a@b.com", Password: "p", Name: "n"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password, name) VALUES (?, ?, ?)")).
		WithArgs(u.Email, u.Password, u.Name).
		WillReturnResult(sqlmock.NewResult(7, 1))

	if err := us.Insert(context.Background(), u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Id != 7 {
		t.Fatalf("expected id 7, got %d", u.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserGet_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = ?")).WithArgs(123).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at"}))

	u, err := us.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u != nil {
		t.Fatalf("expected nil user when no rows, got %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
