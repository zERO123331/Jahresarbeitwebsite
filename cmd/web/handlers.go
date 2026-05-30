package main

import (
	"Jahresarbeitwebsite/internal/models"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	filters := models.Filters{
		PageSize:     2,
		Page:         1,
		Sort:         "id",
		SortSafelist: []string{"id"},
	}
	latestUpdates, err := app.models.Update.GetAll("", filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "home.gohtml", templateData{
		Updates: latestUpdates,
	})
}
