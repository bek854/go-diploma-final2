package api

import (
	"net/http"
	"os"
)

func checkAuth(r *http.Request) bool {
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		return true
	}

	username, password, ok := r.BasicAuth()
	if ok && username == "admin" && password == expectedPassword {
		return true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "Bearer "+expectedPassword {
		return true
	}

	return false
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Планировщик TODO"`)
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
