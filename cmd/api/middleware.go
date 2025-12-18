package main

import (
	"net/http"
	"strings"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/luisfucros/go-events-api-tutorial/internal/configs"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"

)

func (app *application) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    	defer cancel()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			app.unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			app.unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			app.unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		userId := claims["userId"].(float64)

		user, err := app.getUser(ctx, int64(userId))
		if err != nil {
			app.unauthorized(c, "unauthorized access")
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

func (app *application) getUser(ctx context.Context, userId int64) (*store.User, error) {
	if !configs.Envs.REDISEnabled {
		return app.store.Users.Get(ctx, int64(userId))
	}

	user, err := app.cacheStorage.Users.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.Get(ctx, int64(userId))
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
