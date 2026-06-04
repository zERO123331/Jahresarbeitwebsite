package main

import "net/http"

type ErrorPageData struct {
	ErrorID int
	Error   string
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.basicClientError(w, r, http.StatusTooManyRequests)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.stylizedClientError(w, r, http.StatusNotFound, "Not found.")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.stylizedClientError(w, r, http.StatusMethodNotAllowed, "Method not allowed.")
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, page string, form any) {
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, r, http.StatusUnprocessableEntity, page, data)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	app.stylizedClientError(w, r, http.StatusForbidden, "You are not permitted to do that.")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	app.stylizedClientError(w, r, http.StatusBadRequest, "Bad request.")
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	app.stylizedClientError(w, r, http.StatusUnauthorized, "You are not authorized to do that.")
}

func (app *application) failedFilterValidationResponse(w http.ResponseWriter, r *http.Request, page string, filter any) {
	data := app.newTemplateData(r)
	data.Form = filter
	app.render(w, r, http.StatusBadRequest, page, data)
}
