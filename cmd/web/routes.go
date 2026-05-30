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

	dynamic := alice.New()
	protected := dynamic.Append(app.requireAuthentication)

	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/update/view/:id", dynamic.ThenFunc(app.updateView))
	router.Handler(http.MethodGet, "/update/create", protected.ThenFunc(app.updateCreate))
	router.Handler(http.MethodPost, "/update/create", protected.ThenFunc(app.updateCreatePost))
	router.Handler(http.MethodGet, "/update/edit/:id", protected.ThenFunc(app.updateUpdate))
	router.Handler(http.MethodPost, "/update/edit/:id", protected.ThenFunc(app.updateUpdatePost))
	router.Handler(http.MethodGet, "/update/list", dynamic.ThenFunc(app.updateList))

	router.Handler(http.MethodGet, "/user/create", dynamic.ThenFunc(app.userCreate))
	router.Handler(http.MethodPost, "/user/create", dynamic.ThenFunc(app.userCreatePost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/Verify", dynamic.ThenFunc(app.userVerify))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))

	router.Handler(http.MethodGet, "/shop", dynamic.ThenFunc(app.shopPage))
	router.Handler(http.MethodGet, "/shop/view/:id", dynamic.ThenFunc(app.shopEntry))
	router.Handler(http.MethodGet, "/shop/create", protected.ThenFunc(app.shopEntryCreate))
	router.Handler(http.MethodPost, "/shop/create", protected.ThenFunc(app.shopEntryCreatePost))

	// TODO: add favicon and/or find a better way to handle this
	router.HandlerFunc(http.MethodGet, "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	standard := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, commonHeaders, app.sessionManager.LoadAndSave)

	return standard.Then(router)
}
