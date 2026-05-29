package main

import "net/http"

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusTooManyRequests)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusNotFound, "404.gohtml", app.newTemplateData(r))
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusMethodNotAllowed)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.clientError(w, r, http.StatusBadRequest)
}
