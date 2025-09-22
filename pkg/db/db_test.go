package db

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	// Тестируем инициализацию БД
	db, err := Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer db.Close()

	// Тестируем создание задачи
	task := &Task{
		Date:    "20250101",
		Title:   "Test Task",
		Comment: "Test Comment",
		Repeat:  "",
	}

	id, err := db.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if id <= 0 {
		t.Error("Expected positive task ID")
	}

	// Тестируем получение задачи
	retrievedTask, err := db.GetTaskByID("1")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if retrievedTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", retrievedTask.Title)
	}

	// Тестируем получение всех задач
	tasks, err := db.GetAllTasks()
	if err != nil {
		t.Fatalf("Failed to get all tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// Тестируем обновление задачи
	task.ID = "1"
	task.Title = "Updated Task"
	err = db.UpdateTask(task)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Тестируем удаление задачи
	err = db.DeleteTask("1")
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
}

func TestNextDate(t *testing.T) {
	// Этот тест должен быть в пакете api, а не db
	// Оставлен для примера структуры тестов
}
