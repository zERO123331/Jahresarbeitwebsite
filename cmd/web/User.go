package main

import (
	"Jahresarbeitwebsite/internal/models"
	"net/http"
)

type UserCreateForm struct {
	Name        string
	Email       string
	Password    string
	Password2   string
	FieldErrors map[string]string
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserCreateForm{
		FieldErrors: make(map[string]string),
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
		Name:        r.FormValue("username"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Password2:   r.FormValue("password2"),
		FieldErrors: make(map[string]string),
	}

	// TODO: validate form
	// TODO: fix CSS for error messages

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "usercreate.gohtml", data)
		return
	}

	id, err := app.models.User.Insert(&models.User{Name: form.Name, Email: form.Email, Password2: form.Password})
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
