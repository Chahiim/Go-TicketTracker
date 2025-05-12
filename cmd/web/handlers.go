package main

import (
	"fmt"
	"net/http"

	"github.com/chahiim/ticket_tracker/internal/data"
	"github.com/chahiim/ticket_tracker/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Food Ticket System"
	data.HeaderText = "Welcome to the Food Ticket System"
	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) NewTicketHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Create Ticket - Food Ticket System"
	data.HeaderText = "Create a New Food Ticket"

	err := app.render(w, http.StatusOK, "ticket_form.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render Ticket Form page", "template", "ticket_form.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) DisplayTicketHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "View Active Tickets"
	data.HeaderText = "View Active Tickets"
	err := app.render(w, http.StatusOK, "view_tickets.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render journal page", "template", "view_tickets.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createTicket(w http.ResponseWriter,
	r *http.Request) {
	//A. parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	cname := r.PostForm.Get("customerName")
	iname := r.PostForm.Get("itemName")
	quantity := r.PostForm.Get("quantity")

	ticket := &data.Ticket{
		CName:    cname,
		IName:    iname,
		Quantity: quantity,
	}
	//validate data
	v := validator.NewValidator()
	data.ValidateTicket(v, ticket)
	// check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Create Ticket - Food Ticket System"
		data.HeaderText = "Create a New Food Ticket"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"customerName": cname,
			"itemName":     iname,
			"quantity":     quantity,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "home.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.tickets.Insert(ticket)
	if err != nil {
		app.logger.Error("failed to insert ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ticket/success", http.StatusSeeOther)
}

func (app *application) TicketSuccess(w http.ResponseWriter,
	r *http.Request) {
	data := NewTemplateData()
	data.Title = "ticket Submitted"
	data.HeaderText = "Thank You for Your ticket!"
	err := app.render(w, http.StatusOK, "ticket_success.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render ticket success page", "template", "ticket_success.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

/*
func (app *application) DisplayTickets(w http.ResponseWriter, r *http.Request) {

	readTicket :=
		`SELECT *
	FROM ticket
	LIMIT 5;`

	app.db.Query(readTicket)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	defer rows.Close()
}
*/
/*Authentication handlers*/

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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := &data.User{
		Name:           name,
		Email:          email,
		HashedPassword: []byte(password),
	}

	//validate input
	v := validator.NewValidator()
	data.ValidateUser(v, user)
	//check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Login"
		data.HeaderText = "Login Here"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name": name,
			"email": email,
			"password": password,
		}
		err := app.render(w, http.StatusUnprocessableEntity, "login.page.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render login page", "template", "login.page.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.user.Insert(user)
	if err != nil {
		app.logger.Error("failed to insert User", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	app.session.Put(r, "flash","Sign up was successful")
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
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	
	errors_user := make(map[string]string)

	id, err := app.users.Authenticate(email, password)
	if, err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			errors_user["default"] ="Email or Password is Incorrect"
			err := app.render(w, http.StatusOK, "login.page.tmpl", data)
			if err != nil {
				app.logger.Error("failed to render login page", "template", "login.page.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
		return
	}
	app.session.Put(r, "aunthenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out successfully.")
	http.Redirect(w,r, "/user/login", http.StatusSeeOther)
}
