package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Database struct {
	db *sql.DB
}

var DB *Database

func Init(dbFile string) (*Database, error) {
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Используем modernc.org/sqlite (pure Go)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	database := &Database{db: db}
	DB = database
	return database, nil
}

func createTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы tasks: %w", err)
	}
	return nil
}

func (d *Database) AddTask(task *Task) (int64, error) {
	query := `INSERT INTO tasks (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return id, nil
}

func (d *Database) GetTasks(limit int) ([]*Task, error) {
	var query string
	if limit > 0 {
		query = "SELECT id, date, title, comment, repeat FROM tasks ORDER BY date LIMIT ?"
	} else {
		query = "SELECT id, date, title, comment, repeat FROM tasks ORDER BY date"
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = d.db.Query(query, limit)
	} else {
		rows, err = d.db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения задач: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var id int
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования задачи: %w", err)
		}
		task.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по задачам: %w", err)
	}

	return tasks, nil
}

func (d *Database) GetTask(id string) (*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM tasks WHERE id = ?"
	row := d.db.QueryRow(query, id)

	var task Task
	var taskID int
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}

	task.ID = fmt.Sprintf("%d", taskID)
	return &task, nil
}

func (d *Database) UpdateTask(task *Task) error {
	query := "UPDATE tasks SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	result, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func (d *Database) DeleteTask(id string) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func (d *Database) UpdateTaskDate(id, date string) error {
	query := "UPDATE tasks SET date = ? WHERE id = ?"
	result, err := d.db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// CreateTask - алиас для AddTask
func (d *Database) CreateTask(task *Task) (int64, error) {
	return d.AddTask(task)
}

// GetAllTasks - алиас для GetTasks
func (d *Database) GetAllTasks() ([]*Task, error) {
	return d.GetTasks(0)
}

// GetTaskByID - алиас для GetTask
func (d *Database) GetTaskByID(id string) (*Task, error) {
	return d.GetTask(id)
}

// Close - закрытие соединения
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
