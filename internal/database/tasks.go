package database

import (
	"database/sql"
	"fmt"
	"task-manager-go/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db sqlx.DB
}

func NewTaskStore(db sqlx.DB) *TaskStore {
	return &TaskStore{
		db: db,
	}
}

func (t *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
	SELECT id, title, description, completed, updated_at, created_at FROM tasks 
	ORDER BY created_at DESC`

	err := t.db.Select(&tasks, query) // выполняет запрос и автоматически заполняет слайс Task

	if err != nil {
		return nil, fmt.Errorf("Ошибка получения всех задач %w", err)
	}

	return tasks, nil
}

func (t *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	query := `
	SELECT id, title, description, completed, updated_at, created_at 
	FROM tasks WHERE id = $1
    `

	err := t.db.Get(&task, query, id) // выполняет запрос и заполняет одну структуру

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача с id %s не найдена", id)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения по id %w", err)
	}

	return &task, nil
}

func (t *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
	INSERT INTO tasks (title, description, completed, updated_at, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, description, completed, updated_at, created_at
`

	timer := time.Now()

	err := t.db.QueryRowx( // выполняет запрос и возвращает одну строку
		query,
		input.Title,
		input.Description,
		input.Completed,
		timer,
		timer).StructScan(&task)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания задачи %w", err)
	}

	return &task, nil
}
