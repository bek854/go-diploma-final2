package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение задач")

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		log.Printf("Поиск задач: %s", search)
		tasks, err = searchTasks(search, 50)
	} else {
		tasks, err = database.GetTasks(50)
	}

	if err != nil {
		log.Printf("Ошибка при получении задач: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено %d задач", len(tasks))

	if tasks == nil {
		tasks = []*db.Task{}
	}

	response := TasksResp{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)

	log.Println("Ответ отправлен успешно")
}

func searchTasks(search string, limit int) ([]*db.Task, error) {
	allTasks, err := database.GetTasks(0)
	if err != nil {
		return nil, err
	}

	var result []*db.Task

	for _, task := range allTasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(search)) {
			result = append(result, task)
			continue
		}

		if strings.Contains(strings.ToLower(task.Comment), strings.ToLower(search)) {
			result = append(result, task)
			continue
		}

		if isDateMatch(task.Date, search) {
			result = append(result, task)
			continue
		}

		if task.ID == search {
			result = append(result, task)
			continue
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func isDateMatch(taskDate, search string) bool {
	if search == taskDate {
		return true
	}

	if len(search) == 10 && search[2] == '.' && search[5] == '.' {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			searchFormatted := parsedDate.Format("20060102")
			return taskDate == searchFormatted
		}
	}

	return false
}
