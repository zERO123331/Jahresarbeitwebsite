package main

import (
	"Jahresarbeitwebsite/internal/models"
	"Jahresarbeitwebsite/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

type UserCreateForm struct {
	Name                string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Password2           string `form:"password2"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserCreateForm{
		Name: "",
	}
	app.render(w, r, http.StatusOK, "usercreate.gohtml", data)
}

func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {

	var form UserCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckFieldErrors(validator.MinChars(form.Name, 3), "username", fmt.Sprintf("Must be at least %d characters long.", 3))
	form.CheckFieldErrors(validator.NotBlank(form.Name), "username", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Name, 25), "username", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.Matches(form.Name, validator.UserNameRX), "username", "Must be a valid username")
	form.CheckFieldErrors(validator.NotBlank(form.Email), "email", "This field is required.")
	form.CheckFieldErrors(validator.MaxChars(form.Email, 255), "email", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email address")
	form.CheckFieldErrors(validator.NotBlank(form.Password), "password", "This field is required.")
	form.CheckFieldErrors(validator.MinChars(form.Password, 8), "password", fmt.Sprintf("Must be at least %d characters long.", 8))
	form.CheckFieldErrors(validator.MaxChars(form.Password, 72), "password", fmt.Sprintf("Must be at most %d characters long.", 72))
	form.CheckFieldErrors(validator.NotBlank(form.Password2), "password2", "This field is required.")
	form.CheckFieldErrors(form.Password == form.Password2, "password2", "Passwords do not match.")

	// TODO: fix CSS for error messages

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.logger.Error("invalid user create form", "errors", form.FieldErrors)
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

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "userlogin.gohtml", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckFieldErrors(validator.NotBlank(form.Email), "email", "This field is required.")
	form.CheckFieldErrors(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email address")
	form.CheckFieldErrors(validator.MaxChars(form.Email, 255), "email", fmt.Sprintf("Must be at most %d characters long.", 255))
	form.CheckFieldErrors(validator.NotBlank(form.Password), "password", "This field is required.")
	form.CheckFieldErrors(validator.MinChars(form.Password, 8), "password", fmt.Sprintf("Must be at least %d characters long.", 8))
	form.CheckFieldErrors(validator.MaxChars(form.Password, 72), "password", fmt.Sprintf("Must be at most %d characters long.", 72))

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.logger.Error("invalid user login form", "errors", form.FieldErrors)
		app.render(w, r, http.StatusUnprocessableEntity, "userlogin.gohtml", data)
		return
	}

	id, err := app.models.User.Authenticate(form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			form.AddNonFieldError("Email or password is incorrect.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "userlogin.gohtml", data)
		case errors.Is(err, models.ErrUserNotActivated):
			form.AddNonFieldError("Your account has not been activated yet. Please check your email for the activation link.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "userlogin.gohtml", data)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userVerify(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "userverify.gohtml", data)
}
