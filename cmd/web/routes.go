package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("/ticket", app.NewTicketHandler)
	mux.HandleFunc("/ticketsview", app.DisplayTicketHandler)
	mux.HandleFunc("POST /ticket/new", app.createTicket)
	mux.HandleFunc("GET /ticket/display", app.DisplayTickets)
	/*	mux.HandleFunc("GET /ticket/Read", app.ReadTickets)
		mux.HandleFunc("POST /ticket/update", app.UpdateTicket)
		mux.HandleFunc("GET /ticket/delete", app.DeleteTicket)
	*/
	return app.loggingMiddleware(mux)
}
