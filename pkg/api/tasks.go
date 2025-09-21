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

	// Получаем параметр поиска
	search := r.FormValue("search")
	
	var tasks []*db.Task
	var err error

	if search != "" {
		// Поиск задач
		log.Printf("Поиск задач: %s", search)
		tasks, err = searchTasks(search, 50)
	} else {
		// Все задачи
		tasks, err = db.GetTasks(50)
	}

	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено %d задач", len(tasks))

	// Если tasks nil, создаем пустой slice
	if tasks == nil {
		tasks = []*db.Task{}
	}

	// Возвращаем задачи
	response := TasksResp{
		Tasks: tasks,
	}
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
	
	log.Println("Ответ отправлен успешно")
}

// searchTasks ищет задачи по заголовку, комментарию или дате
func searchTasks(search string, limit int) ([]*db.Task, error) {
	allTasks, err := db.GetTasks(0) // 0 = все задачи
	if err != nil {
		return nil, err
	}

	var result []*db.Task
	
	for _, task := range allTasks {
		// Поиск в заголовке
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(search)) {
			result = append(result, task)
			continue
		}
		
		// Поиск в комментарии
		if strings.Contains(strings.ToLower(task.Comment), strings.ToLower(search)) {
			result = append(result, task)
			continue
		}
		
		// Поиск по дате (формат DD.MM.YYYY)
		if isDateMatch(task.Date, search) {
			result = append(result, task)
			continue
		}
		
		// Поиск по ID
		if task.ID == search {
			result = append(result, task)
			continue
		}
	}

	// Ограничиваем количество результатов
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// isDateMatch проверяет совпадение даты в разных форматах
func isDateMatch(taskDate, search string) bool {
	// Прямое совпадение (YYYYMMDD)
	if taskDate == search {
		return true
	}
	
	// Попробуем парсить search как дату в формате DD.MM.YYYY
	if len(search) == 10 && search[2] == '.' && search[5] == '.' {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			// Конвертируем в формат YYYYMMDD для сравнения
			searchFormatted := parsedDate.Format("20060102")
			return taskDate == searchFormatted
		}
	}
	
	return false
}
