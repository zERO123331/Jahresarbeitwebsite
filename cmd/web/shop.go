package main

import (
	"Jahresarbeitwebsite/internal/cdn"
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
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
	Categories          string `form:"categories"`
	validator.Validator `form:"-"`
}

type shopPageInput struct {
	Title      string
	Categories []string
	models.Filters
	validator.Validator
}

func (app *application) shopPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	input := shopPageInput{}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Categories = app.readCSV(qs, "categories", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "price", "quantity", "-id", "-title", "-price", "-quantity"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedFilterValidationResponse(w, r, "shop.gohtml", input)
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
		app.badRequestResponse(w, r)
		return
	}
	entry, err := app.models.Shop.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFoundResponse(w, r)
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

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userID == 0 {
		app.unauthorizedResponse(w, r)
		return
	}

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.stylizedClientError(w, r, http.StatusBadRequest, "invalid form")
		return
	}

	categories := strings.Split(strings.TrimSpace(form.Categories), ",")
	for i, category := range categories {
		category = strings.TrimSpace(category)
		categories[i] = strings.ToLower(category)
	}

	err = r.ParseMultipartForm(1024 * 1024 * 55)
	if err != nil {
		app.logger.Error("failed to parse multipart form", "error", err.Error())
		app.stylizedClientError(w, r, http.StatusBadRequest, "invalid form")
		return
	}
	multipartForm := r.MultipartForm
	if multipartForm == nil {
		app.stylizedClientError(w, r, http.StatusBadRequest, "invalid form")
		return
	}
	fileHeaders := multipartForm.File["images"]

	form.CheckFieldErrors(validator.NotBlank(form.Title), "title", "This field is required.")
	form.CheckFieldErrors(validator.MinChars(form.Title, 10), "title", fmt.Sprintf("Must be at least %d characters long.", 10))
	form.CheckFieldErrors(validator.MaxChars(form.Title, 255), "title", fmt.Sprintf("Must be at most %d characters long.", 255))

	form.CheckFieldErrors(validator.NotBlank(form.Description), "description", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Description, 2000), "description", fmt.Sprintf("Must be at most %d characters long.", 2000))
	form.CheckFieldErrors(validator.MinChars(form.Description, 50), "description", fmt.Sprintf("Must be at least %d characters long.", 50))

	form.CheckFieldErrors(form.Price > 0, "price", "Must be greater than zero.")
	form.CheckFieldErrors(form.Price <= 1000000, "price", "Must be less than or equal to 10,000.00.")

	form.CheckFieldErrors(form.Quantity > 0, "quantity", "Must be greater than zero.")
	form.CheckFieldErrors(form.Quantity <= 1000, "quantity", "Must be less than or equal to 1,000.")

	form.CheckFieldErrors(validator.NotBlank(form.Categories), "categories", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Categories, 255), "categories", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.Unique(categories), "categories", "Each category must be unique.")
	form.CheckFieldErrors(validator.PartEmpty(categories), "categories", "Category list cannot contain empty strings.")
	form.CheckFieldErrors(len(categories) <= 5, "categories", "You can only select up to 5 categories.")

	form.CheckFieldErrors(len(fileHeaders) >= 1, "images", "This field is required.")
	form.CheckFieldErrors(len(fileHeaders) <= 5, "images", "You can only upload up to 5 images.")

	for _, fileHeader := range fileHeaders {
		form.CheckFieldErrors(fileHeader.Size <= 1024*1024*10, "images", "Each image must be smaller than 10MB.")
		form.CheckFieldErrors(fileHeader.Size >= 512*1024, "images", "Each image must be larger than 1MB.")
		contentType := fileHeader.Header.Get("Content-Type")
		form.CheckFieldErrors(cdn.IsAllowedContentType(contentType), "images", "Only JPEG, GIF, WebP and PNG images are allowed.")
	}

	if !form.Valid() {
		app.failedValidationResponse(w, r, "shopEntryCreate.gohtml", form)
		return
	}

	imageURLs, err := app.uploadImages(r, fileHeaders, userID)
	if err != nil {
		app.serverError(w, r, err)
		err = app.deleteImages(r, imageURLs)
		if err != nil {
			panic(err)
		}
		return
	}

	entryID, err := app.models.Shop.Insert(form.Title, form.Description, form.Price, form.Quantity, imageURLs, categories, userID)
	if err != nil {
		app.serverError(w, r, err)
		err = app.deleteImages(r, imageURLs)
		if err != nil {
			panic(err)
		}
		return
	}

	http.Redirect(w, r, "/shop/view/"+strconv.Itoa(entryID), http.StatusSeeOther)
}
