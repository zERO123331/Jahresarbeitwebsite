package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// TODO: add user stuff
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/update/view/:id", dynamic.ThenFunc(app.updateView))
	router.Handler(http.MethodGet, "/update/create", dynamic.ThenFunc(app.updateCreate))
	router.Handler(http.MethodPost, "/update/create", dynamic.ThenFunc(app.updateCreatePost))
	router.Handler(http.MethodGet, "/update/list", dynamic.ThenFunc(app.updates))

	router.Handler(http.MethodGet, "/user/create", dynamic.ThenFunc(app.userCreate))
	router.Handler(http.MethodPost, "/user/create", dynamic.ThenFunc(app.userCreatePost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/Verify", dynamic.ThenFunc(app.userVerify))

	router.Handler(http.MethodGet, "/shop", dynamic.ThenFunc(app.shopPage))

	standard := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, commonHeaders)

	return standard.Then(router)
}
