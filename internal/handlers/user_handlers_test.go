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

//func TestUserHandlers_Update(t *testing.T) {
//	// todo testing
//}
//
//func TestUserHandlers_Delete(t *testing.T) {
//	// todo testing
//}
