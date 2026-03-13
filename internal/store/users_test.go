package store

import (
	"context"
	"errors"
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

func TestUserGet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at"}).
		AddRow(7, "a@b.com", "Alice", "hashed", "2025-01-01")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = ?")).WithArgs(int64(7)).WillReturnRows(rows)

	u, err := us.Get(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil || u.Id != 7 || u.Email != "a@b.com" {
		t.Fatalf("unexpected user: %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserGet_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at"}).
		AddRow("notanint", "a@b.com", "Alice", "hashed", "2025-01-01")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = ?")).WithArgs(int64(7)).WillReturnRows(rows)

	_, err = us.Get(context.Background(), 7)
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserInsert_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}
	u := &User{Email: "a@b.com", Password: "p", Name: "n"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email, password, name) VALUES (?, ?, ?)")).
		WithArgs(u.Email, u.Password, u.Name).
		WillReturnError(errors.New("exec failed"))

	if err := us.Insert(context.Background(), u); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserGetByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at"}).
		AddRow(3, "alice@example.com", "Alice", "hashed", "2025-01-01")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = ?")).
		WithArgs("alice@example.com").WillReturnRows(rows)

	u, err := us.GetByEmail(context.Background(), "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil || u.Id != 3 || u.Email != "alice@example.com" {
		t.Fatalf("unexpected user: %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserGetByEmail_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	us := &UserStore{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = ?")).
		WithArgs("ghost@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at"}))

	u, err := us.GetByEmail(context.Background(), "ghost@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u != nil {
		t.Fatalf("expected nil, got %#v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
