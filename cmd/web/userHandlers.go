package main

import (
	"errors"
	"net/http"

	"github.com/chahiim/ticket_tracker/internal/data"
	"github.com/chahiim/ticket_tracker/internal/validator"
)

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Sign Up "
	data.HeaderText = "Make Your Account Here"
	err := app.render(w, http.StatusOK, "signup.page.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render sign up page", "template", "signup.page.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Create user instance
	users := &data.Users{
		Name:      name,
		Email:     email,
		Activated: true,
	}

	// Validate form data
	v := validator.NewValidator()
	data.ValidateUsers(v, users, password)

	// Show form again if validation failed
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Signup"
		data.HeaderText = "Sign Up Here"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":  name,
			"email": email,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "signup.page.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render signup form", "template", "signup.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Insert user into the database
	err = app.users.Insert(users, password)
	if err != nil {
		app.logger.Error("failed to insert user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.session.Put(r, "flash", "Sign up was successful")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Login"
	data.HeaderText = "Login Here"
	err := app.render(w, http.StatusOK, "login.page.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render login page", "template", "login.page.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse login form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	errors_user := make(map[string]string)
	if email == "" {
		errors_user["email"] = "Email is required"
	}
	if password == "" {
		errors_user["password"] = "Password is required"
	}

	// If validation errors, re-render form
	if len(errors_user) > 0 {
		data := NewTemplateData()
		data.Title = "Login"
		data.HeaderText = "Login"
		data.FormErrors = errors_user
		data.FormData = map[string]string{
			"email": email,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render login form", "template", "login.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Authenticate the user
	user, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, data.ErrInvalidCredentials) {
			data := NewTemplateData()
			data.Title = "Login"
			data.HeaderText = "Login"
			data.FormErrors = map[string]string{
				"generic": "Invalid email or password.",
			}
			data.FormData = map[string]string{
				"email": email,
			}

			err := app.render(w, http.StatusUnauthorized, "login.tmpl", data)
			if err != nil {
				app.logger.Error("failed to render login form", "template", "login.tmpl", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Unknown/internal error
		app.logger.Error("error authenticating user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session
	app.session.Put(r, "user_id", int(user.User_id))
	app.session.Put(r, "authenticatedUserID", true)

	// Redirect to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out successfully.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
