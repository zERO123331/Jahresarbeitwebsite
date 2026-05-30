package main

import (
	"Jahresarbeitwebsite/internal/validator"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	app.render(w, r, http.StatusInternalServerError, "500.gohtml", app.newTemplateData(r))
}

func (app *application) fallbackServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Fallback Server Error", "error", err.Error(), "method", r.Method, "uri", r.URL.RequestURI(), "remote_addr", r.RemoteAddr, "trace", debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %q not found", page)
		app.fallbackServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.fallbackServerError(w, r, err)
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return strings.Split(s, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddFieldError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Get(r.Context(), "authenticatedUserID") != nil
}
