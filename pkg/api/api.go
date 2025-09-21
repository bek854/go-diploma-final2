package api

import (
	"net/http"
)

func Init() {
	// Статические файлы
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// API endpoints с аутентификацией
	http.HandleFunc("/api/nextdate", authMiddleware(nextDateHandler))
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(taskDoneHandler))
}
