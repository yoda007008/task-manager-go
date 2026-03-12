package domain

import "task-manager-go/internal/models"

type UserHandlers interface {
	GetAll() ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Create(input models.CreateTaskInput) (*models.Task, error)
	Update(id int, input models.UpdateTaskInput) (*models.Task, error)
	Delete(id int) error
}
