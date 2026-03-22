package handlers

import (
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

	t.Run("success returning tasks", func(t *testing.T) {
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

	t.Run("returning error", func(t *testing.T) {
		mockRepo.On("GetAll").Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.GetAllTask(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandlers_GetByID(t *testing.T) {
	// todo testing
}

func TestUserHandlers_Create(t *testing.T) {
	// todo testing
}

func TestUserHandlers_Update(t *testing.T) {
	// todo testing
}

func TestUserHandlers_Delete(t *testing.T) {
	// todo testing
}
