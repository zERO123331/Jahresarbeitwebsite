package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) shopPage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string
		Categories []string
		models.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Categories = app.readCSV(qs, "categories", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "price", "quantity"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.FilterErrors)
	}
	entries, err := app.models.Shop.GetAll(input.Title, input.Categories, input.Filters)
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
