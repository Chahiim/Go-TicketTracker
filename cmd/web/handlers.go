package main

import (
	_ "errors"
	"net/http"
	"strconv"

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
		app.logger.Error("failed to render view ticket page", "template", "view_tickets.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
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

	qStr := r.PostForm.Get("quantity")
	qInt, err := strconv.Atoi(qStr)
	if err != nil {
		app.logger.Error("invalid quantity", "value", qStr)
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	ticket := &data.Ticket{
		CName:    cname,
		IName:    iname,
		Quantity: qInt,
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

func (app *application) readTicket(w http.ResponseWriter, r *http.Request) {

	// Call function to get all tickets
	tickets, err := app.tickets.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch tickets", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	flash := app.session.PopString(r, "flash")

	data := NewTemplateData()
	data.TicketList = tickets
	data.Flash = flash

	err = app.render(w, http.StatusOK, "ticket.view.page.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render tickets list", "template", "ticket.view.page.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func (app *application) updateTicket(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ticketid := r.PostForm.Get("ticket_id")
	TicketID, err := strconv.ParseInt(ticketid, 10, 64)
	if err != nil {
		app.logger.Error("invalid ticket id", "value", TicketID)
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	cname := r.PostForm.Get("cname")
	iname := r.PostForm.Get("iname")
	quantity := r.PostForm.Get("quantity")

	qStr := r.PostForm.Get("quantity")
	qInt, err := strconv.Atoi(qStr)
	if err != nil {
		app.logger.Error("invalid quantity", "value", qStr)
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	ticket := &data.Ticket{
		CName:    cname,
		IName:    iname,
		Quantity: qInt,
	}

	v := validator.NewValidator()
	data.ValidateTicket(v, ticket)

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Edit Ticket"
		data.HeaderText = "Edit Ticket"
		data.FormErrors = v.Errors // Store validation errors
		data.FormData = map[string]string{
			"cname":    cname,
			"iname":    iname,
			"quantity": quantity,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "ticket.edit.page.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render edit ticket form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.tickets.Update(ticket)
	if err != nil {
		app.logger.Error("failed to update ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteTicket(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	ticketid := r.FormValue("ticket_id")
	ticketID, err := strconv.ParseInt(ticketid, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	err = app.tickets.Delete(ticketID)
	if err != nil {
		http.Error(w, "Could not delete ticket", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) editTicketForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("ticket_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	ticket, err := app.tickets.GetByID(id)
	if err != nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	data := NewTemplateData()
	data.Title = "Edit Ticket"
	data.HeaderText = "Edit Ticket"
	data.FormData = map[string]string{
		"ticket_id": strconv.FormatInt(ticket.ID, 10),
		"cname":     ticket.CName,
		"iname":     ticket.IName,
		"quantity":  strconv.Itoa(ticket.Quantity),
	}
	app.render(w, http.StatusOK, "ticket.edit.page.tmpl", data)
}
