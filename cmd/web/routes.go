package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("/ticket/new", app.NewTicketHandler)
	mux.HandleFunc("POST /ticket/create", app.createTicket)
	mux.HandleFunc("/ticket/success", app.TicketSuccess)
	mux.HandleFunc("/tickets", app.DisplayTicketHandler)
	mux.HandleFunc("POST /ticket/edit", app.editTicketForm)
	mux.HandleFunc("/ticket/update", app.updateTicket)
	mux.HandleFunc("/ticket/delete", app.deleteTicket)
	mux.HandleFunc("GET /user/signup", app.signupUserForm)
	mux.HandleFunc("POST /user/signup", app.signupUser)
	mux.HandleFunc("GET /user/login", app.loginUserForm)
	mux.HandleFunc("POST /user/login", app.loginUser)
	mux.HandleFunc("POST /user/logout", app.logoutUser)

	return app.loggingMiddleware(mux)
}
