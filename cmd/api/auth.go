package main

import (
	"net/http"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

	var auth loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		app.badRequest(c, err, "invalid format")
		return
	}

	existingUser, err := app.store.Users.GetByEmail(ctx, auth.Email)
	if err != nil {
		app.internalServerError(c, err, "something went wrong")
		return
	}

	if existingUser == nil {
		app.unauthorized(c, "invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		app.unauthorized(c, "invalid email or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"expr":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.config.JWT.Secret))
	if err != nil {
		app.internalServerError(c, err, "error generating token")
		return
	}
	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

func (app *application) registerUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()
	
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		app.badRequest(c, err, "invalid format")
		return
	}

	existingUser, err := app.store.Users.GetByEmail(ctx, register.Email)
	if err != nil {
		app.internalServerError(c, err, "something went wrong")
		return
	}

	if existingUser != nil {
		app.conflict(c, "user already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)

	if err != nil {
		app.internalServerError(c, err, "something went wrong")
		return
	}

	register.Password = string(hashedPassword)
	user := store.User{
		Email:    register.Email,
		Name:     register.Name,
		Password: register.Password,
	}

	err = app.store.Users.Insert(ctx, &user)
	if err != nil {
		app.internalServerError(c, err, "could not create user")
		return
	}
	c.JSON(http.StatusCreated, user)
}
