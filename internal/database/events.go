package database

import "database/sql"

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id           int64   `json:"id"`
	OwnerId      int64   `json:"owner_id" binding:"required"`
	Name         string  `json:"name" binding:"required",min=3`
	Description  string  `json:"description" binding:"required",min=10`
	Date         string  `json:"date" binding:"required",datetime=2006-01-02`
	Location     string  `json:"location" binding:"required",min=3`
}