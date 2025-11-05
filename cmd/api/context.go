package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfucros/go-events-api-tutorial/internal/database"
)

func (app *application) GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exist := c.Get("user")
	if !exist {
		return &database.User{}
	}
	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}
