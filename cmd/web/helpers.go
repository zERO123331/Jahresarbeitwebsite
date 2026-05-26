package main

import (
	"fmt"
	"net/http"
	"time"
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

func GetUpdates() []Update {
	updates := []Update{}
	for i := 1; i <= 25; i++ {
		updates = append(updates, Update{
			Title:       fmt.Sprintf("Test Update %d", i),
			Author:      "Test Author",
			Body:        "Test Body",
			Created:     time.Now(),
			LastUpdated: time.Now(),
			ID:          i,
		})
	}
	return updates
}

func GetLatestUpdate() Update {
	return Update{
		Title:       "Test Update",
		Author:      "Test Author",
		Body:        "Test Body",
		Created:     time.Now(),
		LastUpdated: time.Now(),
		ID:          25,
	}
}

func GetUpdateByID(id int) (Update, error) {
	if id > 25 {
		return Update{}, fmt.Errorf("update %d not found", id)
	}
	return Update{
		Title:       fmt.Sprintf("Test Update %d", id),
		Author:      "Test Author",
		Body:        "Test Body",
		Created:     time.Now(),
		LastUpdated: time.Now(),
		ID:          id,
	}, nil
}
