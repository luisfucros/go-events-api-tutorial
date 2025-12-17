package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	app.logger.Info("setting up routes")

	g := gin.Default()
	g.Use(app.RateLimiter())

	v1 := g.Group("/api/v1")
	{
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)

		v1.GET("/attendees/:id/events", app.getEventsByAttendee)
		v1.GET("/events/:id/attendees/", app.getAttendeesForEvent)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)

		authGroup.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json"))(c)
	})

	return g
}
