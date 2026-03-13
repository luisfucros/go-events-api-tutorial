package store

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEventInsert_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}
	evt := &Event{OwnerId: 10, Name: "Test", Description: "A long description", Date: "2025-12-25", Location: "Here"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO events (owner_id, name, description, date, location) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(evt.OwnerId, evt.Name, evt.Description, evt.Date, evt.Location).
		WillReturnResult(sqlmock.NewResult(42, 1))

	if err := es.Insert(context.Background(), evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evt.Id != 42 {
		t.Fatalf("expected id 42, got %d", evt.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventInsert_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}
	evt := &Event{OwnerId: 10, Name: "Test", Description: "A long description", Date: "2025-12-25", Location: "Here"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO events (owner_id, name, description, date, location) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(evt.OwnerId, evt.Name, evt.Description, evt.Date, evt.Location).
		WillReturnError(errors.New("exec failed"))

	if err := es.Insert(context.Background(), evt); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).
		AddRow(1, 2, "one", "desc one", "2025-01-01", "loc1").
		AddRow(2, 3, "two", "desc two", "2025-01-02", "loc2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events")).WillReturnRows(rows)

	list, err := es.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 events, got %d", len(list))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGet_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events WHERE id = ?")).WithArgs(123).WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}))

	e, err := es.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if e != nil {
		t.Fatalf("expected nil event when no rows, got %#v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGet_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	// Return a row with a wrong type to cause scan error
	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).AddRow("notint", 1, "n", "d", "2025-01-01", "l")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events WHERE id = ?")).WithArgs(123).WillReturnRows(rows)

	_, err = es.Get(context.Background(), 123)
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventUpdate_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}
	evt := &Event{Id: 1, Name: "n", Description: "d", Date: "2025-01-01", Location: "l"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE events SET name = ?, description = ?, date = ?, location = ? WHERE id = ?")).
		WithArgs(evt.Name, evt.Description, evt.Date, evt.Location, evt.Id).
		WillReturnError(errors.New("exec failed"))

	if err := es.Update(context.Background(), evt); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventDelete_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM events WHERE id = ?")).WithArgs(5).WillReturnError(errors.New("exec failed"))

	if err := es.Delete(context.Background(), 5); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).
		AddRow(1, 2, "Test Event", "A description here", "2025-06-01", "NYC")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events WHERE id = ?")).WithArgs(int64(1)).WillReturnRows(rows)

	e, err := es.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil || e.Id != 1 || e.Name != "Test Event" {
		t.Fatalf("unexpected event: %#v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGetAll_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events")).WillReturnError(errors.New("query failed"))

	_, err = es.GetAll(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventGetAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "date", "location"}).
		AddRow("notint", 2, "n", "d", "2025-01-01", "l")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM events")).WillReturnRows(rows)

	_, err = es.GetAll(context.Background())
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}
	evt := &Event{Id: 1, Name: "Updated", Description: "New description", Date: "2025-07-01", Location: "LA"}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE events SET name = ?, description = ?, date = ?, location = ? WHERE id = ?")).
		WithArgs(evt.Name, evt.Description, evt.Date, evt.Location, evt.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := es.Update(context.Background(), evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestEventDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	es := &EventStore{DB: db}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM events WHERE id = ?")).WithArgs(int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := es.Delete(context.Background(), 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
