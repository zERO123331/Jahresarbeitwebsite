package main

import (
	"fmt"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	app.render(w, r, http.StatusInternalServerError, "500.gohtml", app.newTemplateData(r))
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %q not found", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{}
}
