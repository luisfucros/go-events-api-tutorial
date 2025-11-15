package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Events interface {
		Get(context.Context, int64) (*Event, error)
		GetAll(context.Context) ([]*Event, error)
		Insert(context.Context, *Event) error
		Delete(context.Context, int64) error
		Update(context.Context, *Event) error
	}
	Users interface {
		Get(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Insert(context.Context, *User) error
	}
	Attendees interface {
		GetByEventAndAttendee(context.Context, int64, int64) (*Attendee, error)
		GetAttendeesByEvent(context.Context, int64) ([]*User, error)
		Insert(context.Context, *Attendee) (*Attendee, error)
		Delete(context.Context, int64, int64) error
		GetEventsByAttendee(context.Context, int64) ([]*Event, error)
		
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Events:    &EventStore{db},
		Users:    &UserStore{db},
		Attendees: &AttendeeStore{db},
	}
}