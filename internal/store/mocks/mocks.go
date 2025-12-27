package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
)

// MockEvents is a mock implementation of the Events store
type MockEvents struct{ mock.Mock }

func (m *MockEvents) Get(ctx context.Context, id int64) (*store.Event, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*store.Event), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockEvents) GetAll(ctx context.Context) ([]*store.Event, error) {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.([]*store.Event), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockEvents) Insert(ctx context.Context, e *store.Event) error { args := m.Called(ctx, e); return args.Error(0) }
func (m *MockEvents) Update(ctx context.Context, e *store.Event) error { args := m.Called(ctx, e); return args.Error(0) }
func (m *MockEvents) Delete(ctx context.Context, id int64) error { args := m.Called(ctx, id); return args.Error(0) }

// MockUsers is a mock implementation of the Users store
type MockUsers struct{ mock.Mock }

func (m *MockUsers) Get(ctx context.Context, id int64) (*store.User, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*store.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUsers) GetByEmail(ctx context.Context, email string) (*store.User, error) {
	args := m.Called(ctx, email)
	if v := args.Get(0); v != nil {
		return v.(*store.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUsers) Insert(ctx context.Context, u *store.User) error { args := m.Called(ctx, u); return args.Error(0) }

// MockAttendees is a mock implementation of the Attendees store
type MockAttendees struct{ mock.Mock }

func (m *MockAttendees) GetByEventAndAttendee(ctx context.Context, eventId, userId int64) (*store.Attendee, error) {
	args := m.Called(ctx, eventId, userId)
	if v := args.Get(0); v != nil {
		return v.(*store.Attendee), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAttendees) GetAttendeesByEvent(ctx context.Context, eventId int64) ([]*store.User, error) {
	args := m.Called(ctx, eventId)
	if v := args.Get(0); v != nil {
		return v.([]*store.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAttendees) Insert(ctx context.Context, a *store.Attendee) (*store.Attendee, error) {
	args := m.Called(ctx, a)
	if v := args.Get(0); v != nil {
		return v.(*store.Attendee), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAttendees) Delete(ctx context.Context, userId, eventId int64) error { args := m.Called(ctx, userId, eventId); return args.Error(0) }
func (m *MockAttendees) GetEventsByAttendee(ctx context.Context, id int64) ([]*store.Event, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.([]*store.Event), args.Error(1)
	}
	return nil, args.Error(1)
}

// Cache mocks

type MockCacheEvents struct{ mock.Mock }
func (m *MockCacheEvents) Get(ctx context.Context, id int64) (*store.Event, error) { args := m.Called(ctx, id); if v := args.Get(0); v != nil { return v.(*store.Event), args.Error(1) }; return nil, args.Error(1) }
func (m *MockCacheEvents) Set(ctx context.Context, e *store.Event) error { args := m.Called(ctx, e); return args.Error(0) }
func (m *MockCacheEvents) Delete(ctx context.Context, id int64) { m.Called(ctx, id) }
func (m *MockCacheEvents) GetAll(ctx context.Context) ([]store.Event, error) { args := m.Called(ctx); if v := args.Get(0); v != nil { return v.([]store.Event), args.Error(1) }; return nil, args.Error(1) }
func (m *MockCacheEvents) SetAll(ctx context.Context, events []store.Event) error { args := m.Called(ctx, events); return args.Error(0) }
func (m *MockCacheEvents) DeleteAll(ctx context.Context) { m.Called(ctx) }

type MockCacheUsers struct{ mock.Mock }
func (m *MockCacheUsers) Get(ctx context.Context, id int64) (*store.User, error) { args := m.Called(ctx, id); if v := args.Get(0); v != nil { return v.(*store.User), args.Error(1) }; return nil, args.Error(1) }
func (m *MockCacheUsers) Set(ctx context.Context, u *store.User) error { args := m.Called(ctx, u); return args.Error(0) }
func (m *MockCacheUsers) Delete(ctx context.Context, id int64) { m.Called(ctx, id) }

type MockCacheAttendees struct{ mock.Mock }
func (m *MockCacheAttendees) GetAttendeesByEvent(ctx context.Context, eventID int64) ([]store.User, error) { args := m.Called(ctx, eventID); if v := args.Get(0); v != nil { return v.([]store.User), args.Error(1) }; return nil, args.Error(1) }
func (m *MockCacheAttendees) SetAttendeesByEvent(ctx context.Context, eventID int64, users []store.User) error { args := m.Called(ctx, eventID, users); return args.Error(0) }
func (m *MockCacheAttendees) DeleteAttendeesByEvent(ctx context.Context, eventID int64) { m.Called(ctx, eventID) }
func (m *MockCacheAttendees) GetEventsByAttendee(ctx context.Context, userID int64) ([]store.Event, error) { args := m.Called(ctx, userID); if v := args.Get(0); v != nil { return v.([]store.Event), args.Error(1) }; return nil, args.Error(1) }
func (m *MockCacheAttendees) SetEventsByAttendee(ctx context.Context, userID int64, events []store.Event) error { args := m.Called(ctx, userID, events); return args.Error(0) }
func (m *MockCacheAttendees) DeleteEventsByAttendee(ctx context.Context, userID int64) { m.Called(ctx, userID) }
