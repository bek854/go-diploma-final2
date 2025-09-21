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
	
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Инициализация API обработчиков
	api.Init()

	// Жёстко задаем порт 9092
	port := "9092"

	// Запуск веб-сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
