package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на отметку выполнения задачи")
	
	if r.Method != http.MethodPost {
		log.Printf("Неподдерживаемый метод: %s", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		log.Println("Не указан ID задачи")
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	log.Printf("Отметка выполнения задачи ID: %s", id)

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("Ошибка получения задачи: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Если задача не повторяющаяся, удаляем ее
	if task.Repeat == "" {
		log.Printf("Удаление одноразовой задачи ID: %s", id)
		err = db.DeleteTask(id)
		if err != nil {
			log.Printf("Ошибка удаления задачи: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Одноразовая задача ID: %s удалена", id)
	} else {
		// Для повторяющейся задачи вычисляем следующую дату
		log.Printf("Перенос повторяющейся задачи ID: %s, правило: %s", id, task.Repeat)
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка вычисления следующей даты: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			log.Printf("Ошибка обновления даты задачи: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Повторяющаяся задача ID: %s перенесена на %s", id, nextDate)
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{})
	
	log.Printf("Задача ID: %s успешно обработана", id)
}
