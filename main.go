package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	// Инициализация базы данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	database, err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Сохраняем ссылку на БД в api пакете
	api.SetDB(database)

	// Инициализация обработчиков API
	api.Init()

	// Порт
	port := "8082"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	// Запуск веб-сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
