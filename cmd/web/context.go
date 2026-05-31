package main

import (
	"Jahresarbeitwebsite/internal/models"
	"context"
	"net/http"
)

type contextKey string

const (
	isAuthenticatedContextKey contextKey = "isAuthenticated"
	userContextKey            contextKey = "user"
)

func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userContextKey, user))
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
