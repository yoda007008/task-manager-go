package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"task-manager-go/internal/database"
	"task-manager-go/internal/models"
)

type TaskHandlers struct {
	store *database.TaskStore
}

func NewTaskHandler(store *database.TaskStore) *TaskHandlers {
	return &TaskHandlers{
		store: store,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	// поле payload - данные, которые нужно отправлять в формате json
	// w http.ResponseWriter - используется для отправки ответа
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload) // создаем новый Encoder, который пишет в ResponseWriter, Encode преобразует поле в jsonи отправляет ответ
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// endpoint GET /tasks
func (t *TaskHandlers) GetAllTask(w http.ResponseWriter, r *http.Request) {
	task, err := t.store.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// endpoint GET /tasks/{id} <- этот айди надо извлечь из url пути
func (t *TaskHandlers) GetTask(w http.ResponseWriter, r *http.Request) {
	// тут мы разбиваем путь по айди
	// разбиваем путь r.URL.Path по "/" и берем последний элемент id -> GET /tasks/{id}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	task, err := t.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// endpoint POST /task -> новая задача
func (t *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	err := json.NewDecoder(r.Body).Decode(&input) // декодируем тело запроса и вставляем в структуру input
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок к задаче обязателен")
		return
	}

	task, err := t.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Ошибка создания задачи")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// endpoint PUT /tasks/{id}
func (t *TaskHandlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	var input models.UpdateTaskInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	// проверяем, что заголовок не пустой
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок не может быть пустым")
		return
	}

	task, err := t.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, "Обновление задачи произошло некорректно")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// endpoint DELETE /task/{id}
func (t *TaskHandlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	err = t.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, "Удаление задачи произошло некорректно")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Задача с ID %d успешно удалена", id),
	})
}
