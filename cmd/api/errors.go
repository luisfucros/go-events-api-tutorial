package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) internalServerError(c *gin.Context, err error, message string) {
	app.logger.Errorw(
		"internal server error",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"error", err.Error(),
	)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func (app *application) badRequest(c *gin.Context, err error, message string) {
	app.logger.Warnw(
		"bad request",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"error", err.Error(),
	)

	c.JSON(http.StatusBadRequest, gin.H{
		"error": message,
	})
}

func (app *application) unauthorized(c *gin.Context, message string) {
	app.logger.Warnw(
		"unauthorized",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": message,
	})
}

func (app *application) forbidden(c *gin.Context, message string) {
	app.logger.Warnw(
		"forbidden",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	c.JSON(http.StatusForbidden, gin.H{
		"error": message,
	})

}

func (app *application) notFound(c *gin.Context, message string) {
	app.logger.Warnw(
		"not found",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	c.JSON(http.StatusNotFound, gin.H{
		"error": message,
	})
}

func (app *application) conflict(c *gin.Context, message string) {
	app.logger.Warnw(
		"conflict",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	c.JSON(http.StatusConflict, gin.H{
		"error": message,
	})
}

func (app *application) rateLimitExceeded(c *gin.Context, retryAfter string) {
	app.logger.Warnw(
		"rate limit exceeded",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	c.Header("Retry-After", retryAfter)
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": "rate limit exceeded",
	})
}
