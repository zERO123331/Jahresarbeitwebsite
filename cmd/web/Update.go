package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UpdateCreateForm struct {
	Title string `form:"title"`
	Body  string `form:"body"`
	validator.Validator
}

func (app *application) updateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
	data.Update = update
	app.render(w, r, http.StatusOK, "viewupdate.gohtml", data)

}

// TODO: sometimes a server error occurs when transfering to this page
func (app *application) updateCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UpdateCreateForm{
		Title: "",
		Body:  "",
	}
	app.render(w, r, http.StatusOK, "updatecreate.gohtml", data)
}

func (app *application) updateCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := UpdateCreateForm{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}

	form.CheckFieldErrors(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Title, 255), "title", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.NotBlank(form.Body), "body", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Body, 10000), "body", fmt.Sprintf("Must be at most %d characters long.", 10000))
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "updatecreate.gohtml", data)
		return
	}

	update := &models.Update{
		Title:  form.Title,
		Body:   form.Body,
		UserID: 15, // TODO: implement user
	}

	id, err := app.models.Update.Insert(update.UserID, update.Title, update.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/update/view/"+strconv.Itoa(id), http.StatusSeeOther)
	app.logger.Info("update created", "id", id)
}

func (app *application) updates(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
	data.Updates = updates
	app.render(w, r, http.StatusOK, "updates.gohtml", data)
}
