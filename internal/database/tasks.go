package database

import (
	"database/sql"
	"fmt"
	"task-manager-go/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
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

func (t *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := t.GetByID(id)

	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	var updatedTask models.Task

	task.UpdatedAt = time.Now()

	query := `UPDATE tasks SET title = $1, description = $2, completed = $3, updated_at = $4
	WHERE id = $5
	RETURNING id, title, description, created_at, updated_at`

	err = t.db.QueryRowx(
		query,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Completed,
		task.UpdatedAt,
		id,
	).StructScan(&updatedTask)

	if err != nil {
		return nil, fmt.Errorf("ошибка обновления задачи", err)
	}

	return &updatedTask, nil
}

func (t *TaskStore) Delete(id int) error {

	query := `
	DELETE FROM tasks
	WHERE id = $1
	`

	result, err := t.db.Exec(query, id) // функция Exeс выполняет запрос, не возвращая ни одной строки

	if err != nil {
		return fmt.Errorf("ошибка проверки результата удаления задачи", err)
	}

	// RowsAffected возвращает кол-во затронутых (удаленных) строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления", err)
	}

	if rowsAffected == 0 {
		fmt.Errorf("задача с id %d не найдена", id)
	}

	return nil
}
