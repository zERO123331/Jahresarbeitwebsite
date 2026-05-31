package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type shopEntryCreateForm struct {
	Title               string `form:"title"`
	Description         string `form:"description"`
	Price               int    `form:"price"`
	Quantity            int    `form:"quantity"`
	ImageURLs           string `form:"image_urls"`
	Categories          string `form:"categories"`
	validator.Validator `form:"-"`
}

func (app *application) shopPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
	input.Filters.SortSafelist = []string{"id", "title", "price", "quantity", "-id", "-title", "-price", "-quantity"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, "shop.gohtml", input)
		return
	}
	entries, err := app.models.Shop.GetAll(input.Title, input.Categories, input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.ShopEntries = entries
	app.render(w, r, http.StatusOK, "shop.gohtml", data)
}

func (app *application) shopEntry(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id < 1 {
		app.stylizedClientError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid shop entry id: %s", r.PathValue("id")))
		return
	}
	entry, err := app.models.Shop.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.stylizedClientError(w, r, http.StatusNotFound, "shop entry not found")
			return
		}
		app.serverError(w, r, err)
		return
	}
	data.ShopEntry = entry
	app.render(w, r, http.StatusOK, "shopEntryView.gohtml", data)
}

func (app *application) shopEntryCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = shopEntryCreateForm{
		Title: "",
	}
	app.render(w, r, http.StatusOK, "shopEntryCreate.gohtml", data)
}

func (app *application) shopEntryCreatePost(w http.ResponseWriter, r *http.Request) {
	var form shopEntryCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.basicClientError(w, r, http.StatusBadRequest)
		return
	}

	categories := strings.Split(strings.TrimSpace(form.Categories), ",")
	for i, category := range categories {
		category = strings.TrimSpace(category)
		categories[i] = category
	}

	form.CheckFieldErrors(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Title, 255), "title", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.NotBlank(form.Description), "description", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Description, 2000), "description", fmt.Sprintf("Must be at most %d characters long.", 2000))
	form.CheckFieldErrors(validator.MinChars(form.Description, 50), "description", fmt.Sprintf("Must be at least %d characters long.", 50))
	form.CheckFieldErrors(form.Price > 0, "price", "This field is required.")
	form.CheckFieldErrors(form.Quantity > 0, "quantity", "This field is required.")
	form.CheckFieldErrors(validator.NotBlank(form.Categories), "categories", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Categories, 255), "categories", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.Unique(categories), "categories", "Each category must be unique.")
	//form.CheckFieldErrors(validator.NotBlank(form.ImageURLs), "image_urls", "This field is required.")
	//form.CheckFieldErrors(validator.MaxChars(form.ImageURLs, 255), "image_urls", fmt.Sprintf("Must be at most %d characters long.", 255))
	// images will come later but at this point idk how to effectively handle them, without storing local files or using a separate CDN

	if !form.Valid() {
		app.failedValidationResponse(w, r, "shopEntryCreate.gohtml", form)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	imageURLs := []string{"/static/img/placeholder.svg"}

	entryID, err := app.models.Shop.Insert(form.Title, form.Description, form.Price, form.Quantity, imageURLs, categories, userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/shop/view/"+strconv.Itoa(entryID), http.StatusSeeOther)
}
