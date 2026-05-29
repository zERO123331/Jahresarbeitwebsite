package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	latestUpdates, err := app.models.Update.GetLatest(2)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "home.gohtml", templateData{
		Updates: latestUpdates,
	})
}
