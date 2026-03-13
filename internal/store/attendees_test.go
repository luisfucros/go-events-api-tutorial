package store

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAttendeeInsert_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}
	a := &Attendee{EventId: 1, UserId: 2}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO attendees (event_id, user_id) VALUES (?, ?)")).
		WithArgs(a.EventId, a.UserId).
		WillReturnResult(sqlmock.NewResult(99, 1))

	res, err := as.Insert(context.Background(), a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Id != 99 {
		t.Fatalf("expected id 99, got %d", res.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetByEventAndAttendee_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, event_id FROM attendees WHERE event_id = ? AND user_id = ?")).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "event_id"}))

	res, err := as.GetByEventAndAttendee(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil, got %#v", res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetAttendeesByEvent_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow("badid", "n", "e")
	mock.ExpectQuery(regexp.QuoteMeta("\n\t\tSELECT u.id, u.name, u.email\n\t\tFROM users u\n\t\tJOIN attendees a ON u.id = a.user_id\n\t\twhere a.event_id = ?\n\t")).WithArgs(10).WillReturnRows(rows)

	_, err = as.GetAttendeesByEvent(context.Background(), 10)
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeDelete_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM attendees WHERE user_id = ? AND event_id=?")).WithArgs(1, 2).WillReturnError(errors.New("exec fail"))

	if err := as.Delete(context.Background(), 1, 2); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetEventsByAttendee_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).
		AddRow(1, 10, "a", "d", "2025-01-01", "l")

	mock.ExpectQuery(regexp.QuoteMeta("\n\t\tSELECT e.id, e.owner_id, e.name, e.description, e.date, e.location\n\t\tFROM events e\n\t\tJOIN attendees a ON e.id = a.event_id\n\t\tWHERE a.user_id = ?\n\t")).WithArgs(5).WillReturnRows(rows)

	list, err := as.GetEventsByAttendee(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 1 || list[0].Id != 1 {
		t.Fatalf("unexpected list: %#v", list)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeInsert_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}
	a := &Attendee{EventId: 1, UserId: 2}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO attendees (event_id, user_id) VALUES (?, ?)")).
		WithArgs(a.EventId, a.UserId).
		WillReturnError(errors.New("exec failed"))

	_, err = as.Insert(context.Background(), a)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetByEventAndAttendee_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "event_id"}).AddRow(10, 2, 1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, event_id FROM attendees WHERE event_id = ? AND user_id = ?")).
		WithArgs(1, 2).WillReturnRows(rows)

	res, err := as.GetByEventAndAttendee(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.Id != 10 {
		t.Fatalf("unexpected result: %#v", res)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetAttendeesByEvent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Alice", "alice@example.com").
		AddRow(2, "Bob", "bob@example.com")
	mock.ExpectQuery(regexp.QuoteMeta("\n\t\tSELECT u.id, u.name, u.email\n\t\tFROM users u\n\t\tJOIN attendees a ON u.id = a.user_id\n\t\twhere a.event_id = ?\n\t")).WithArgs(5).WillReturnRows(rows)

	users, err := as.GetAttendeesByEvent(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 || users[0].Name != "Alice" || users[1].Name != "Bob" {
		t.Fatalf("unexpected users: %#v", users)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM attendees WHERE user_id = ? AND event_id=?")).
		WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := as.Delete(context.Background(), 1, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestAttendeeGetEventsByAttendee_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	as := &AttendeeStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).
		AddRow("notint", 10, "a", "d", "2025-01-01", "l")
	mock.ExpectQuery(regexp.QuoteMeta("\n\t\tSELECT e.id, e.owner_id, e.name, e.description, e.date, e.location\n\t\tFROM events e\n\t\tJOIN attendees a ON e.id = a.event_id\n\t\tWHERE a.user_id = ?\n\t")).WithArgs(5).WillReturnRows(rows)

	_, err = as.GetEventsByAttendee(context.Background(), 5)
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
