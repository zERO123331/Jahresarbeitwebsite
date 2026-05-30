package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) updateView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.serverError(w, r, fmt.Errorf("invalid update id: %s", params.ByName("id")))
		return
	}
	update, err := app.models.Update.GetByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "viewupdate.gohtml", templateData{
		Update: update,
	})

}

// TODO: implement update Create
func (app *application) updateCreate(w http.ResponseWriter, r *http.Request) {

}

// TODO: implement update Create Post
func (app *application) updateCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func (app *application) updates(w http.ResponseWriter, r *http.Request) {
	var input struct {
		models.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "updated_at"}
	models.ValidateFilters(v, input.Filters)

	updates, err := app.models.Update.GetAll("", input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "updates.gohtml", templateData{
		Updates: updates,
	})
}
