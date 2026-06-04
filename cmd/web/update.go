package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type updateCreateForm struct {
	Title               string `form:"title"`
	Body                string `form:"body"`
	validator.Validator `form:"-"`
}

type updateUpdateForm struct {
	Title               string `form:"title"`
	Body                string `form:"body"`
	ID                  int    `form:"id"`
	Version             int    `form:"version"`
	validator.Validator `form:"-"`
}

func (app *application) updateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.badRequestResponse(w, r)
		return
	}
	update, err := app.models.Update.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data.Update = update
	app.render(w, r, http.StatusOK, "viewupdate.gohtml", data)

}

// TODO: sometimes a server error occurs when transfering to this page
func (app *application) updateCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = updateCreateForm{
		Title: "",
		Body:  "",
	}
	app.render(w, r, http.StatusOK, "updatecreate.gohtml", data)
}

func (app *application) updateCreatePost(w http.ResponseWriter, r *http.Request) {

	var form updateCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.stylizedClientError(w, r, http.StatusBadRequest, "invalid form")
		return
	}

	form.CheckFieldErrors(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Title, 255), "title", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.NotBlank(form.Body), "body", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Body, 10000), "body", fmt.Sprintf("Must be at most %d characters long.", 10000))
	if !form.Valid() {
		app.failedValidationResponse(w, r, "updatecreate.gohtml", form)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	update := &models.Update{
		Title:  form.Title,
		Body:   form.Body,
		UserID: userID,
	}

	id, err := app.models.Update.Insert(update.UserID, update.Title, update.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/update/view/"+strconv.Itoa(id), http.StatusSeeOther)
	app.logger.Info("update created", "id", id)
}

func (app *application) updateList(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	var input struct {
		models.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "updated_at", "-id", "-title", "-updated_at"}
	models.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, "updates.gohtml", input)
		return
	}

	updates, err := app.models.Update.GetAll("", input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Updates = updates
	app.render(w, r, http.StatusOK, "updates.gohtml", data)
}

func (app *application) updateUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.badRequestResponse(w, r)
		return
	}
	update, err := app.models.Update.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data.Update = update
	data.Form = updateUpdateForm{
		Title: "",
	}
	app.render(w, r, http.StatusOK, "updateupdate.gohtml", data)
}

func (app *application) updateUpdatePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.badRequestResponse(w, r)
		return
	}

	var form updateUpdateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.stylizedClientError(w, r, http.StatusBadRequest, "invalid form")
		return
	}

	form.CheckFieldErrors(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Title, 255), "title", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.NotBlank(form.Body), "body", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Body, 10000), "body", fmt.Sprintf("Must be at most %d characters long.", 10000))
	if !form.Valid() {
		app.failedValidationResponse(w, r, "updateupdate.gohtml", form)
		return
	}

	err = app.models.Update.Update(id, form.Title, form.Body, form.Version)
	if err != nil {
		if errors.Is(err, models.ErrEditConflict) {
			app.sessionManager.Put(r.Context(), "flash", "Tried to edit an outdated version of the update. Please try again.")
			app.logger.Info("update edit conflict", "id", id, "version", form.Version)
			http.Redirect(w, r, "/update/view/"+strconv.Itoa(id), http.StatusSeeOther)
			return
		}
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/update/view/"+strconv.Itoa(id), http.StatusSeeOther)
}
