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
