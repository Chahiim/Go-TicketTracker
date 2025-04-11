package main

import (
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

func (app *application) ticketHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) ticketViewerHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.ticket.Insert(ticket)
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
