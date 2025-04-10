package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("/feedback", app.feedbackHandler)
	mux.HandleFunc("/journal", app.journalHandler)
	mux.HandleFunc("/todo", app.todoHandler)
	mux.HandleFunc("POST /feedback/new", app.createFeedback)
	mux.HandleFunc("GET /feedback/success", app.feedbackSuccess)
	mux.HandleFunc("POST /journal/new", app.createJournal)
	mux.HandleFunc("GET /journal/Success", app.journalSuccess)
	mux.HandleFunc("POST /todo/new", app.createTodo)
	mux.HandleFunc("GET /todo/success", app.todoSuccess)

	return app.loggingMiddleware(mux)
}
