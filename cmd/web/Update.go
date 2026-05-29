package main

import (
	"Jahresarbeitwebsite/internal/models"
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
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// TODO: implement update Create Post
func (app *application) updateCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

// TODO: rework updates to use pages and search stuff
func (app *application) updates(w http.ResponseWriter, r *http.Request) {
	updates, err := app.models.Update.GetAll("", models.Filters{})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "updates.gohtml", templateData{
		Updates: updates,
	})
}
