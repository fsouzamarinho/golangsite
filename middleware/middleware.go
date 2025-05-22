package middleware

import (
	"andravirtual/handlers"
	"net/http"
)

func RequireLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, logged := handlers.UsuarioLogado(r)
		if !logged {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}