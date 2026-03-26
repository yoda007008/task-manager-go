package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task-manager-go/internal/models"
	"task-manager-go/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHandlers_GetAllTask(t *testing.T) {
	mockRepo := mocks.NewTaskService(t)

	handler := NewTaskHandler(mockRepo)

	t.Run("success returning tasks, status 200", func(t *testing.T) {
		expectedTasks := []models.Task{
			{ID: 1, Title: "Task 1", Completed: false},
			{ID: 2, Title: "Task 2", Completed: true},
		}

		mockRepo.On("GetAll").Return(expectedTasks, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.GetAllTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response []models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)

		require.NoError(t, err)
		assert.Equal(t, expectedTasks, response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returning internal server error status code 500", func(t *testing.T) {
		mockRepo.On("GetAll").Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.GetAllTask(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandlers_GetByID(t *testing.T) {
	mockRepo := mocks.NewTaskService(t)

	handler := NewTaskHandler(mockRepo)

	t.Run("success returning by id status 200", func(t *testing.T) {
		expectedTask := &models.Task{
			ID:    1,
			Title: "Test Test",
		}

		mockRepo.On("GetByID", 1).Return(expectedTask, nil).Once()

		resp := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()

		handler.GetTask(w, resp)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)

		require.NoError(t, err)
		assert.Equal(t, expectedTask.ID, response.ID)
		assert.Equal(t, expectedTask.Title, response.Title)

		mockRepo.AssertExpectations(t)
	})

	t.Run("not task in database returning status 400", func(t *testing.T) {
		resp := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.GetTask(w, resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandlers_Create(t *testing.T) {
	mockRepo := mocks.NewTaskService(t)

	handler := NewTaskHandler(mockRepo)

	t.Run("success create task, status 200", func(t *testing.T) {
		input := models.CreateTaskInput{
			Title:       "First task",
			Description: "Run this program success",
			Completed:   true,
		}

		expectedTask := &models.Task{
			ID:          1,
			Title:       input.Title,
			Description: input.Description,
			Completed:   input.Completed,
		}

		mockRepo.On("Create", input).Return(expectedTask, nil).Once()

		body, err := json.Marshal(input)
		require.NoError(t, err)

		resp := httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateTask(w, resp)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.Task
		err = json.Unmarshal(w.Body.Bytes(), &response)

		require.NoError(t, err)
		assert.Equal(t, expectedTask.ID, response.ID)
		assert.Equal(t, expectedTask.Title, response.Title)
		assert.Equal(t, expectedTask.Description, expectedTask.Description)
		assert.Equal(t, expectedTask.Completed, expectedTask.Completed)

		mockRepo.AssertExpectations(t)
	})

	t.Run("not success response, status 400", func(t *testing.T) {
		resp := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.CreateTask(w, resp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandlers_Update(t *testing.T) {
	mockRepo := mocks.NewTaskService(t)

	handler := NewTaskHandler(mockRepo)

	newTitle := "New title"
	newDescription := "New description"
	newCompleted := true

	t.Run("success update task, status 200", func(t *testing.T) {
		input := models.UpdateTaskInput{
			Title:       &newTitle,
			Description: &newDescription,
			Completed:   &newCompleted,
		}

		expectedTask := &models.Task{
			ID:          1,
			Title:       newTitle,
			Description: newDescription,
			Completed:   newCompleted,
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		mockRepo.On("Update", 1, input).Return(expectedTask, nil).Once()

		resp := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))

		resp.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.UpdateTask(w, resp)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.Task
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, expectedTask.ID, response.ID)
		assert.Equal(t, expectedTask.Title, response.Title)
		assert.Equal(t, expectedTask.Description, response.Description)
		assert.Equal(t, expectedTask.Completed, response.Completed)

		mockRepo.AssertExpectations(t)
	})

	t.Run("not success update task, status 400", func(t *testing.T) {
		input := models.UpdateTaskInput{
			Title: &newTitle,
		}

		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/tasks/abc", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse map[string]string
		json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Contains(t, errResponse["error"], "Некорректный ID")

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("empty title after trim, returns 400", func(t *testing.T) {
		emptyTitle := "   "
		input := models.UpdateTaskInput{
			Title: &emptyTitle,
		}

		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResponse map[string]string
		json.Unmarshal(w.Body.Bytes(), &errResponse)
		assert.Contains(t, errResponse["error"], "Заголовок не может быть пустым")

		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestUserHandlers_Delete(t *testing.T) {
	mockRepo := mocks.NewTaskService(t)

	handler := NewTaskHandler(mockRepo)

	t.Run("success delete task task, status 200", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response["message"], "успешно удалена")
		assert.Contains(t, response["message"], "1")

		mockRepo.AssertExpectations(t)
	})

	t.Run("not found task for delete, status 500", func(t *testing.T) {
		mockRepo.On("Delete", 999).Return(errors.New("задача с id 999 не найдена")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "не найдена")

		mockRepo.AssertExpectations(t)
	})

	t.Run("negative ID, status 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/-5", nil)
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Отрицательный ID не принимается")

		mockRepo.AssertNotCalled(t, "Delete")
	})
}
