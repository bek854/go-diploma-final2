package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var (
	tasks     []*Task
	tasksMu   sync.Mutex
	nextID    int = 1
	taskFile     = "tasks.json"
)

func Init(dbFile string) error {
	tasks = []*Task{}

	// Пытаемся загрузить задачи из файла, если он существует
	if _, err := os.Stat(taskFile); err == nil {
		data, err := os.ReadFile(taskFile)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла задач: %v", err)
		}

		if len(data) > 0 {
			err = json.Unmarshal(data, &tasks)
			if err != nil {
				return fmt.Errorf("ошибка парсинга задач: %v", err)
			}

			// Находим максимальный ID
			maxID := 0
			for _, task := range tasks {
				if id := atoi(task.ID); id > maxID {
					maxID = id
				}
			}
			nextID = maxID + 1
		}
	}

	return nil
}

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации задач: %v", err)
	}

	err = os.WriteFile(taskFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи задач: %v", err)
	}

	return nil
}

func atoi(s string) int {
	var n int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		}
	}
	return n
}

func AddTask(task *Task) (int64, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	task.ID = fmt.Sprintf("%d", nextID)
	nextID++
	tasks = append(tasks, task)

	err := saveTasks()
	if err != nil {
		// Откатываем изменения при ошибке
		tasks = tasks[:len(tasks)-1]
		nextID--
		return 0, err
	}

	return int64(atoi(task.ID)), nil
}

func GetTasks(limit int) ([]*Task, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	// Если задач нет, возвращаем пустой слайс
	if len(tasks) == 0 {
		return []*Task{}, nil
	}

	if limit <= 0 || limit > len(tasks) {
		limit = len(tasks)
	}

	result := make([]*Task, limit)
	copy(result, tasks[:limit])

	return result, nil
}

func GetTask(id string) (*Task, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, fmt.Errorf("задача не найдена")
}

func UpdateTask(updatedTask *Task) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			tasks[i] = updatedTask
			return saveTasks()
		}
	}

	return fmt.Errorf("задача не найдена")
}

func DeleteTask(id string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTasks()
		}
	}

	return fmt.Errorf("задача не найдена")
}

func UpdateTaskDate(id, date string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	for _, task := range tasks {
		if task.ID == id {
			task.Date = date
			return saveTasks()
		}
	}

	return fmt.Errorf("задача не найдена")
}
