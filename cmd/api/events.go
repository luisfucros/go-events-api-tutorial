package main

import (
	"net/http"
	"strconv"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
	"github.com/luisfucros/go-events-api-tutorial/internal/configs"
)

// createEvent creates a new event
//
// @Summary Create a new event
// @Description Creates a new event for the logged-in user
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body database.Event true "Event to create"
// @Success 201 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events [post]
func (app *application) createEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var event store.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		app.badRequest(c, err, "bad request")
		return
	}
	

	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id
	err := app.store.Events.Insert(ctx, &event)

	if err != nil {
		app.internalServerError(c, err, "failed to create event")
		return
	}

	// invalidate cache
    app.cacheInvalidateEvent(ctx, event.Id)

	c.JSON(http.StatusCreated, event)
}

// getEvents return all events
//
// @Summary Get all events
// @Description Returns all events in the db
// @Tags Events
// @Produce json
// @Success 200 {object} []database.Event
// @Failure 500 {object} map[string]string
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	// ---- CACHE READ ----
    if cached, ok := app.cacheGetAllEvents(ctx); ok {
        c.JSON(http.StatusOK, cached)
        return
    }

	events, err := app.store.Events.GetAll(ctx)

	if err != nil {
		app.internalServerError(c, err, "failed to retrieve events")
		return
	}

	eventList := make([]store.Event, len(events))
	for i, e := range events {
		eventList[i] = *e
	}

	// ---- CACHE SET ----
    app.cacheSetAllEvents(ctx, eventList)

	c.JSON(http.StatusOK, events)
}

// getEvent returns a single event by ID
//
// @Summary Get an event by ID
// @Description Returns details for a specific event
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.badRequest(c, err, "invalid event id")
		return
	}

	// --- CACHE READ ---
    if cached, ok := app.cacheGetEvent(ctx, id); ok {
        c.JSON(http.StatusOK, cached)
        return
    }

	event, err := app.store.Events.Get(ctx, id)

	if err != nil {
		app.internalServerError(c, err, "failed to retrieve event")
		return
	}

	if event == nil {
		app.notFound(c, "event not found")
		return
	}

	// --- CACHE SET ---
    app.cacheSetEvent(ctx, event)

	c.JSON(http.StatusOK, event)
}


// updateEvent updates an existing event
//
// @Summary Update an existing event
// @Description Updates an event owned by the logged-in user
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param event body database.Event true "Updated event data"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.badRequest(c, err, "invalid event id")
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.store.Events.Get(ctx, id)

	if err != nil {
		app.internalServerError(c, err, "something went wrong")
		return
	}

	if existingEvent == nil {
		app.notFound(c, "event not found")
		return
	}

	if existingEvent.OwnerId != user.Id {
		app.forbidden(c, "not authorized to update this event")
		return
	}

	updatedEvent := &store.Event{}

	if err := c.ShouldBindBodyWithJSON(updatedEvent); err != nil {
		app.badRequest(c, err, "something went wrong")
		return
	}

	updatedEvent.Id = id

	if err := app.store.Events.Update(ctx, updatedEvent); err != nil {
		app.internalServerError(c, err, "failed to update event")
		return
	}

	// ---- CACHE INVALIDATIONS ----
    app.cacheInvalidateEvent(ctx, id)

	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent deletes an event
//
// @Summary Delete an event
// @Description Deletes an event owned by the logged-in user
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.badRequest(c, err, "invalid event id")
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.store.Events.Get(ctx, id)

	if err != nil {
		app.internalServerError(c, err, "failed to retrieve event")
		return
	}
	if existingEvent == nil {
		app.notFound(c, "event not found")
		return
	}

	if existingEvent.OwnerId != user.Id {
		app.forbidden(c, "")
		return
	}

	if err := app.store.Events.Delete(ctx, id); err != nil {
		app.internalServerError(c, err, "failed to delete event")
		return
	}

	// ---- CACHE INVALIDATIONS ----
    app.cacheInvalidateEvent(ctx, id)

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	eventId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid event id")
		return
	}

	idParam = c.Param("userId")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid user id")
		return
	}

	event, err := app.store.Events.Get(ctx, eventId)
	if err != nil {
		app.internalServerError(c, err, "failed to retrieve event")
		return
	}
	if event == nil {
		app.notFound(c, "event not found")
		return
	}

	userToAdd, err := app.store.Users.Get(ctx, userId)
	if err != nil {
		app.internalServerError(c, err, "failed to retrieve user")
		return
	}
	if userToAdd == nil {
		app.notFound(c, "user not found")
		return
	}

	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		app.forbidden(c, "not authorized to add an attendee")
		return
	}

	existingAttendee, err := app.store.Attendees.GetByEventAndAttendee(ctx, event.Id, userToAdd.Id)
	if err != nil {
		app.internalServerError(c, err, "failed to retrieve attendee")
		return
	}
	if existingAttendee != nil {
		app.conflict(c, "Attendee already exists")
		return
	}

	attendee := store.Attendee{
		EventId: event.Id,
		UserId: userToAdd.Id,
	}

	_, err = app.store.Attendees.Insert(ctx, &attendee)
	if err != nil {
		app.internalServerError(c, err, "failed to add attendee to event")
		return
	}

	// ---- CACHE INVALIDATION ----
	app.cacheInvalidateAttendeesByEvent(ctx, event.Id)
	app.cacheInvalidateEventsByAttendee(ctx, userToAdd.Id)

	c.JSON(http.StatusCreated, attendee)
}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	eventId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid event id")
		return
	}

	// ---- CACHE READ ----
	if cached, ok := app.cacheGetAttendeesByEvent(ctx, eventId); ok {
		c.JSON(http.StatusOK, cached)
		return
	}

	users, err := app.store.Attendees.GetAttendeesByEvent(ctx, eventId)
	if err != nil {
		app.internalServerError(c, err, "failed to retrieve attendees for event")
		return
	}

	userList := make([]store.User, len(users))
	for i, u := range users {
		userList[i] = *u
	}

	// ---- CACHE SET ----
	app.cacheSetAttendeesByEvent(ctx, eventId, userList)

	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	eventId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid event id")
		return
	}

	idParam = c.Param("userId")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid user id")
		return
	}

	user := app.GetUserFromContext(c)
	event, err := app.store.Events.Get(ctx, eventId)

	if err != nil {
		app.internalServerError(c, err, "failed to retrieve event")
		return
	}

	if event == nil {
		app.notFound(c, "event not found")
		return
	}

	if event.OwnerId != user.Id {
		app.forbidden(c, "not authorized to delete an attendee from event")
		return
	}

	err = app.store.Attendees.Delete(ctx, userId, eventId)
	if err != nil {
		app.internalServerError(c, err, "failed to delete attendee")
		return
	}

	// ---- CACHE INVALIDATION ----
	app.cacheInvalidateAttendeesByEvent(ctx, eventId)
	app.cacheInvalidateEventsByAttendee(ctx, userId)

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventsByAttendee(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequest(c, err, "invalid attendee id")
		return
	}

	// ---- CACHE READ ----
	if cached, ok := app.cacheGetEventsByAttendee(ctx, id); ok {
		c.JSON(http.StatusOK, cached)
		return
	}

	events, err := app.store.Attendees.GetEventsByAttendee(ctx, id)
	if err != nil {
		app.internalServerError(c, err, "failed to get events")
		return
	}

	eventList := make([]store.Event, len(events))
	for i, e := range events {
		eventList[i] = *e
	}

	// ---- CACHE SET ----
	app.cacheSetEventsByAttendee(ctx, id, eventList)

	c.JSON(http.StatusOK, events)
}

func (app *application) cacheGetAllEvents(ctx context.Context) ([]store.Event, bool) {
    if !configs.Envs.REDISEnabled {
        return nil, false
    }
    list, err := app.cacheStorage.Events.GetAll(ctx)
    if err != nil || list == nil {
        return nil, false
    }
    return list, true
}

func (app *application) cacheSetAllEvents(ctx context.Context, events []store.Event) {
    if !configs.Envs.REDISEnabled {
        return
    }
    _ = app.cacheStorage.Events.SetAll(ctx, events)
}

func (app *application) cacheGetEvent(ctx context.Context, id int64) (*store.Event, bool) {
    if !configs.Envs.REDISEnabled {
        return nil, false
    }
    cached, err := app.cacheStorage.Events.Get(ctx, id)
    if err != nil || cached == nil {
        return nil, false
    }
    return cached, true
}

func (app *application) cacheSetEvent(ctx context.Context, event *store.Event) {
    if !configs.Envs.REDISEnabled {
        return
    }
    _ = app.cacheStorage.Events.Set(ctx, event)
}

func (app *application) cacheInvalidateEvent(ctx context.Context, id int64) {
    if !configs.Envs.REDISEnabled {
        return
    }
    app.cacheStorage.Events.Delete(ctx, id)
    app.cacheStorage.Events.DeleteAll(ctx)
}

func (app *application) cacheGetAttendeesByEvent(
	ctx context.Context,
	eventID int64,
) ([]store.User, bool) {

	if !configs.Envs.REDISEnabled {
		return nil, false
	}

	users, err := app.cacheStorage.Attendees.GetAttendeesByEvent(ctx, eventID)
	if err != nil || users == nil {
		return nil, false
	}

	return users, true
}

func (app *application) cacheSetAttendeesByEvent(
	ctx context.Context,
	eventID int64,
	users []store.User,
) {
	if !configs.Envs.REDISEnabled {
		return
	}
	_ = app.cacheStorage.Attendees.SetAttendeesByEvent(ctx, eventID, users)
}

func (app *application) cacheInvalidateAttendeesByEvent(
	ctx context.Context,
	eventID int64,
) {
	if !configs.Envs.REDISEnabled {
		return
	}
	app.cacheStorage.Attendees.DeleteAttendeesByEvent(ctx, eventID)
}

func (app *application) cacheGetEventsByAttendee(
	ctx context.Context,
	userID int64,
) ([]store.Event, bool) {

	if !configs.Envs.REDISEnabled {
		return nil, false
	}

	events, err := app.cacheStorage.Attendees.GetEventsByAttendee(ctx, userID)
	if err != nil || events == nil {
		return nil, false
	}

	return events, true
}

func (app *application) cacheSetEventsByAttendee(
	ctx context.Context,
	userID int64,
	events []store.Event,
) {
	if !configs.Envs.REDISEnabled {
		return
	}
	_ = app.cacheStorage.Attendees.SetEventsByAttendee(ctx, userID, events)
}

func (app *application) cacheInvalidateEventsByAttendee(
	ctx context.Context,
	userID int64,
) {
	if !configs.Envs.REDISEnabled {
		return
	}
	app.cacheStorage.Attendees.DeleteEventsByAttendee(ctx, userID)
}
