package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	task, err := database.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(task)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на добавление задачи")

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Получена задача: %+v", task)

	if task.Title == "" {
		log.Println("Заголовок задачи не указан")
		http.Error(w, "Заголовок задачи не указан", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	} else {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			log.Printf("Неверный формат даты: %s", task.Date)
			http.Error(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
	}

	taskDate, _ := time.Parse(dateFormat, task.Date)
	if taskDate.Before(now) && task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка при вычислении следующей даты: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	}

	id, err := database.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в БД: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Задача добавлена с идентификатором: %d", id)

	response := map[string]interface{}{"id": fmt.Sprintf("%d", id)}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, "Заголовок задачи не указан", http.StatusBadRequest)
		return
	}

	if task.Date != "" {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
	}

	err = database.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	err := database.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
