package main

import (
	"Jahresarbeitwebsite/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	latestUpdate := GetLatestUpdate()
	app.render(w, r, http.StatusOK, "home.gohtml", templateData{
		LatestUpdate: latestUpdate,
	})
}
func (app *application) updateView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.serverError(w, r, fmt.Errorf("invalid update id: %s", params.ByName("id")))
		return
	}
	update, err := GetUpdateByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.render(w, r, http.StatusOK, "viewupdate.gohtml", templateData{
		SpecificUpdate: update,
	})

}
func (app *application) updateCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}
func (app *application) updateCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func (app *application) updates(w http.ResponseWriter, r *http.Request) {
	updates := GetUpdates()
	app.render(w, r, http.StatusOK, "updates.gohtml", templateData{
		Updates: updates,
	})
}

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
		SpecificEntry: entry,
	})
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "usercreate.gohtml", data)
}

func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	name := r.FormValue("username")
	if name == "" {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	if email == "" {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password2")
	if password != passwordConfirm {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	id, err := app.models.User.Insert(&models.User{Name: name, Email: email, Password2: password})
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/user/Verify/", http.StatusSeeOther)
	app.logger.Info("user created", "id", id)
}

func (app *application) userVerify(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "userverify.gohtml", data)
}
