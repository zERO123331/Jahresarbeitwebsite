package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"fmt"
	"net/http"
)

type UserCreateForm struct {
	Name      string
	Email     string
	Password  string
	Password2 string
	validator.Validator
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserCreateForm{
		Name: "",
	}
	app.render(w, r, http.StatusOK, "usercreate.gohtml", data)
}

func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form := UserCreateForm{
		Name:      r.FormValue("username"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		Password2: r.FormValue("password2"),
	}

	form.Check(validator.NotBlank(form.Name), "username", "This field is required.")
	form.Check(validator.MaxChars(form.Name, 255), "username", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.Check(validator.NotBlank(form.Email), "email", "This field is required.")
	form.Check(validator.MaxChars(form.Email, 255), "email", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email address")
	form.Check(validator.NotBlank(form.Password), "password", "This field is required.")
	form.Check(validator.MinChars(form.Password, 8), "password", fmt.Sprintf("Must be at least %d characters long.", 8))
	form.Check(validator.MaxChars(form.Password, 72), "password", fmt.Sprintf("Must be at most %d characters long.", 72))
	form.Check(validator.NotBlank(form.Password2), "password2", "This field is required.")
	form.Check(form.Password == form.Password2, "password2", "Passwords do not match.")

	// TODO: fix CSS for error messages

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.logger.Error("invalid user create form", "errors", form.Errors)
		app.render(w, r, http.StatusUnprocessableEntity, "usercreate.gohtml", data)
		return
	}

	user := &models.User{Name: form.Name, Email: form.Email}
	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	id, err := app.models.User.Insert(user)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		app.logger.Error("failed to create user", "error", err.Error())
		return
	}
	http.Redirect(w, r, "/user/Verify/", http.StatusSeeOther)
	app.logger.Info("user created", "id", id)
}

func (app *application) userVerify(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "userverify.gohtml", data)
}
