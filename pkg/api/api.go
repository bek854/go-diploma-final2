package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

var database *db.Database

func SetDB(db *db.Database) {
	database = db
}

func Init() {
	// Статические файлы
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Конечные точки API с аутентификацией
	http.HandleFunc("/api/nextdate", authMiddleware(nextDateHandler))
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(taskDoneHandler))
}
