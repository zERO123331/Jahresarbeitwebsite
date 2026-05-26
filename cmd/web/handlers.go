package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	latestUpdate := GetLatestUpdate()
	app.render(w, r, http.StatusOK, "home.gohtml", templateData{
		Update: latestUpdate,
	})
}
