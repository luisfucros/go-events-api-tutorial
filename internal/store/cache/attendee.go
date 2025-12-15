package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
)

type AttendeeStore struct {
	rdb *redis.Client
}

const (
	EventAttendeesTTL   = 1 * time.Minute
	AttendeeEventsTTL   = 1 * time.Minute
)

func (s *AttendeeStore) GetAttendeesByEvent(
	ctx context.Context,
	eventID int64,
) ([]store.User, error) {

	key := fmt.Sprintf("event:%d:attendees", eventID)

	data, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var users []store.User
	if err := json.Unmarshal([]byte(data), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *AttendeeStore) SetAttendeesByEvent(
	ctx context.Context,
	eventID int64,
	users []store.User,
) error {

	key := fmt.Sprintf("event:%d:attendees", eventID)

	payload, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, key, payload, EventAttendeesTTL).Err()
}

func (s *AttendeeStore) DeleteAttendeesByEvent(
	ctx context.Context,
	eventID int64,
) {
	key := fmt.Sprintf("event:%d:attendees", eventID)
	s.rdb.Del(ctx, key)
}

func (s *AttendeeStore) GetEventsByAttendee(
	ctx context.Context,
	userID int64,
) ([]store.Event, error) {

	key := fmt.Sprintf("attendee:%d:events", userID)

	data, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var events []store.Event
	if err := json.Unmarshal([]byte(data), &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *AttendeeStore) SetEventsByAttendee(
	ctx context.Context,
	userID int64,
	events []store.Event,
) error {

	key := fmt.Sprintf("attendee:%d:events", userID)

	payload, err := json.Marshal(events)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, key, payload, AttendeeEventsTTL).Err()
}

func (s *AttendeeStore) DeleteEventsByAttendee(
	ctx context.Context,
	userID int64,
) {
	key := fmt.Sprintf("attendee:%d:events", userID)
	s.rdb.Del(ctx, key)
}
