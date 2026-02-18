package handlers

import "task-manager-go/internal/database"

type TaskHandlers struct {
	store database.TaskStore
}

func respondWithJSON() {
	// todo вспомогательная функция для отправки хэндлера
}

func respondWithError() {
	// todo
}
