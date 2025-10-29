package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	v1.POST("/events", app.createEvent)
	v1.GET("/events", app.getAllEvents)
	v1.GET("/event/:id", app.getEvent)
	v1.PUT("/event/:id", app.updateEvent)
	v1.DELETE("/event/:id", app.deleteEvent)
	v1.POST("/event/:id/attendees/:userId", app.addAttendeeToEvent)
	v1.DELETE("/event/:id/attendees/:userId", app.deleteAttendeeFromEvent)
	v1.GET("/attendees/:id/events", app.getEventsByAttendee)
	v1.GET("/event/:id/attendees/", app.getAttendeesForEvent)

	v1.POST("/auth/register", app.registerUser)

	return g
}
