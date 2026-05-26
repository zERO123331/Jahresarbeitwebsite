package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/update/view/:id", app.updateView)
	router.HandlerFunc(http.MethodGet, "/update/create", app.updateCreate)
	router.HandlerFunc(http.MethodPost, "/update/create", app.updateCreatePost)
	router.HandlerFunc(http.MethodGet, "/update/list", app.updates)

	router.HandlerFunc(http.MethodGet, "/user/create", app.userCreate)
	router.HandlerFunc(http.MethodPost, "/user/create", app.userCreatePost)
	router.HandlerFunc(http.MethodGet, "/user/Verify", app.userVerify)

	router.HandlerFunc(http.MethodGet, "/shop", app.shopPage)
	return rateLimit(commonHeaders(router))
}
