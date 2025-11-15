package store

import (
	"context"
	"errors"
	"database/sql"
)

type AttendeeStore struct {
	DB *sql.DB
}

type Attendee struct {
	Id         int64   `json:"id"`
	UserId     int64   `json:"user_id"`
	EventId    int64   `json:"event_id"`
}

func(m *AttendeeStore) Insert(ctx context.Context, attendee *Attendee) (*Attendee, error) {
	query := "INSERT INTO attendees (event_id, user_id) VALUES (?, ?)"

	result, err := m.DB.ExecContext(ctx, query,
		attendee.EventId, attendee.UserId,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	attendee.Id = id
	return attendee, nil
}

func(m *AttendeeStore) GetByEventAndAttendee(ctx context.Context, eventId, userId int64) (*Attendee, error) {
	query := `SELECT id, user_id, event_id FROM attendees WHERE event_id = ? AND user_id = ?`

	row := m.DB.QueryRowContext(ctx, query, eventId, userId)

	var attendee Attendee
	err := row.Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeeStore) GetAttendeesByEvent(ctx context.Context, eventId int64) ([]*User, error) {
	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		where a.event_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name, &user.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *AttendeeStore) Delete(ctx context.Context, userId, eventId int64) error {	
	query := "DELETE FROM attendees WHERE user_id = ? AND event_id=?"
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}

	return nil
}

func (m *AttendeeStore) GetEventsByAttendee(ctx context.Context, attendeeId int64) ([]*Event, error) {	
	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, attendeeId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}