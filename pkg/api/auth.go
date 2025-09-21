package api

import (
	"net/http"
	"os"
)

// Базовая проверка аутентификации
func checkAuth(r *http.Request) bool {
	// Если пароль не установлен, аутентификация не требуется
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		return true
	}
	
	// Проверяем базовую аутентификацию
	username, password, ok := r.BasicAuth()
	if ok && username == "admin" && password == expectedPassword {
		return true
	}
	
	// Проверяем заголовок Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "Bearer "+expectedPassword {
		return true
	}
	
	return false
}

// Middleware для аутентификации
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="TODO Planner"`)
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
