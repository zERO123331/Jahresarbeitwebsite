package main

import (
	"Jahresarbeitwebsite/internal/validator"
	"bytes"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	app.render(w, r, http.StatusInternalServerError, "500.gohtml", app.newTemplateData(r))
}

func (app *application) fallbackServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Fallback Server Error", "error", err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) basicClientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) stylizedClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	data := app.newTemplateData(r)
	data.Error = ErrorPageData{
		ErrorID: status,
		Error:   message,
	}

	app.render(w, r, status, "error.gohtml", data)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %q not found", page)
		app.fallbackServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.fallbackServerError(w, r, err)
		return
	}
	w.WriteHeader(status)
	if _, err = buf.WriteTo(w); err != nil {
		app.logger.Error("failed to write response", "error", err.Error())
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return strings.Split(s, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddFieldError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)

	if err != nil {
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) getFormImage(r *http.Request) (image.Image, error) {
	err := r.ParseMultipartForm(1024 * 1024)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("image")

	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (app *application) uploadImages(r *http.Request, fileHeaders []*multipart.FileHeader, userID int) ([]string, error) {
	imageURLS := make([]string, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			return imageURLS, err
		}
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		extension := filepath.Ext(fileHeader.Filename)
		objectKey := fmt.Sprintf("shop/%d/%s%s", userID, uuid.NewString(), extension)
		imageURL, err := app.cdn.Upload(r.Context(), file, fileHeader.Size, objectKey, contentType)
		if err != nil {
			return imageURLS, err
		}
		imageURLS = append(imageURLS, imageURL)
	}
	return imageURLS, nil
}

func (app *application) deleteImages(r *http.Request, imageUrls []string) error {
	for _, imageUrl := range imageUrls {
		err := app.cdn.Delete(r.Context(), imageUrl)
		if err != nil {
			return err
		}
	}
	return nil
}
