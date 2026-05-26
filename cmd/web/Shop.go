package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) shopPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	entries, err := app.models.Shop.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "shop.gohtml", templateData{
		ShopEntries: entries,
	})
}

func (app *application) shopEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.serverError(w, r, fmt.Errorf("invalid shop entry id: %s", r.PathValue("id")))
		return
	}
	entry, err := app.models.Shop.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "shopentry.gohtml", templateData{
		ShopEntry: entry,
	})
}
