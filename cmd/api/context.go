package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
)

func (app *application) GetUserFromContext(c *gin.Context) *store.User {
	contextUser, exist := c.Get("user")
	if !exist {
		return &store.User{}
	}
	user, ok := contextUser.(*store.User)
	if !ok {
		return &store.User{}
	}

	return user
}
