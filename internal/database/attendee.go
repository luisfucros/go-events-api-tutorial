package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id        int64   `json:"id"`
	UserId    int64   `json:"user_id"`
	EventId    int64   `json:"event_id"`      
}