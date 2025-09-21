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
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
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
		log.Printf("Ошибка разбора JSON: %v", err)
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Получена задача: %+v", task)

	// Проверяем обязательное поле
	if task.Title == "" {
		log.Println("Заголовок задачи не указан")
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Устанавливаем дату по умолчанию если не указана
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	} else {
		// Проверяем формат даты
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			log.Printf("Неверный формат даты: %s", task.Date)
			http.Error(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
	}

	// Если дата в прошлом и есть правило повторения, вычисляем следующую дату
	taskTime, _ := time.Parse(dateFormat, task.Date)
	if taskTime.Before(now) && task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка вычисления следующей даты: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка добавления задачи в БД: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Задача добавлена с ID: %d", id)

	// Возвращаем ID созданной задачи
	response := map[string]interface{}{"id": fmt.Sprintf("%d", id)}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if task.ID == "" {
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверяем формат даты если указана
	if task.Date != "" {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
	}

	// Обновляем задачу в БД
	err = db.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
