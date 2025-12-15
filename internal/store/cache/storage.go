package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64)
	}
	Events interface {
        Get(context.Context, int64) (*store.Event, error)
        Set(context.Context, *store.Event) error
        Delete(context.Context, int64)
		GetAll(ctx context.Context) ([]store.Event, error)
		SetAll(ctx context.Context, events []store.Event) error
		DeleteAll(ctx context.Context)
    }
	Attendees interface {
		// event -> users
		GetAttendeesByEvent(ctx context.Context, eventID int64) ([]store.User, error)
		SetAttendeesByEvent(ctx context.Context, eventID int64, users []store.User) error
		DeleteAttendeesByEvent(ctx context.Context, eventID int64)

		// user -> events
		GetEventsByAttendee(ctx context.Context, userID int64) ([]store.Event, error)
		SetEventsByAttendee(ctx context.Context, userID int64, events []store.Event) error
		DeleteEventsByAttendee(ctx context.Context, userID int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
		Events: &EventStore{rdb: rbd},
		Attendees: &AttendeeStore{rdb: rbd},
	}
}