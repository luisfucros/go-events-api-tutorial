package store

import (
	"context"
	"database/sql"
)

type EventStore struct {
	DB *sql.DB
}

type Event struct {
	Id          int64  `json:"id"`
	OwnerId     int64  `json:"ownerId"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
}

func (m *EventStore) Insert(ctx context.Context, event *Event) error {
	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES (?, ?, ?, ?, ?)"

	result, err := m.DB.ExecContext(ctx, query,
		event.OwnerId, event.Name, event.Description, event.Date, event.Location,
	)
	
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	
	if err != nil {
		return err
	}

	event.Id = id
	return nil
}

func (m *EventStore) GetAll(ctx context.Context) ([]*Event, error) {
	query := "SELECT * FROM events"

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

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

func (m *EventStore) Get(ctx context.Context, id int64) (*Event, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := "SELECT * FROM events WHERE id = ?"

	var event Event
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func (m *EventStore) Update(ctx context.Context, event *Event) error {
	query := "UPDATE events SET name = ?, description = ?, date = ?, location = ? WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)

	if err != nil {
		return err
	}

	return nil
}

func (m *EventStore) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM events WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
