package domain

import (
	"net/http"
	"task-manager-go/internal/models"
)

type UserHandlers interface {
	GetAllTask(w http.ResponseWriter, _ *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type TaskService interface {
	GetAll() ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Create(input models.CreateTaskInput) (*models.Task, error)
	Update(id int, input models.UpdateTaskInput) (*models.Task, error)
	Delete(id int) error
}
