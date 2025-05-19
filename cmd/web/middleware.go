package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
		app.logger.Info("Request processed")
	})
	return fn

}

/*
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Use the isAuthenticated
		if !app.isAuthenticated(r) {
			app.logger.Warn("Authentication required", "uri", r.URL.RequestURI()) // Log attempt

			// Redirect the user to the login page.
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
*/
// noSurf middleware
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",                  // Available across the entire site
		Secure:   true,                 // Requires HTTPS
		SameSite: http.SameSiteLaxMode, // Standard SameSite setting
	})

	return csrfHandler
}
