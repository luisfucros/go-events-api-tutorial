package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
)

type EventStore struct {
	rdb *redis.Client
}

const (
    EventListKey = "events-all"
    EventListExp = time.Minute
	EventExpTime = time.Minute
)

func (s *EventStore) Get(ctx context.Context, eventID int64) (*store.Event, error) {
	cacheKey := fmt.Sprintf("event-%d", eventID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var event store.Event
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *EventStore) Set(ctx context.Context, event *store.Event) error {
	cacheKey := fmt.Sprintf("event-%d", event.Id)

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, jsonData, EventExpTime).Err()
}

func (s *EventStore) Delete(ctx context.Context, eventID int64) {
	cacheKey := fmt.Sprintf("event-%d", eventID)
	s.rdb.Del(ctx, cacheKey)
}

// ---- list caching ----

func (s *EventStore) GetAll(ctx context.Context) ([]store.Event, error) {
	data, err := s.rdb.Get(ctx, EventListKey).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, err
	}

	var events []store.Event
	if err := json.Unmarshal([]byte(data), &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventStore) SetAll(ctx context.Context, events []store.Event) error {
	jsonData, err := json.Marshal(events)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, EventListKey, jsonData, EventListExp).Err()
}

func (s *EventStore) DeleteAll(ctx context.Context) {
	s.rdb.Del(ctx, EventListKey)
}