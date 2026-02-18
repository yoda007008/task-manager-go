package handlers

import (
	"encoding/json"
	"net/http"
	"task-manager-go/internal/database"
)

type TaskHandlers struct {
	store database.TaskStore
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	// поле payload - данные, которые нужно отправлять в формате json
	// w http.ResponseWriter - используется для отправки ответа
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload) // создаем новый Encoder, который пишет в ResponseWriter, Encode преобразует поле в jsonи отправляет ответ
}

func respondWithError() {
	// todo
}

func (t *TaskHandlers) GetAllTask(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (t *TaskHandlers) GetTask(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (t *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (t *TaskHandlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (t *TaskHandlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// todo
}
