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
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
		Events: &EventStore{rdb: rbd},
	}
}