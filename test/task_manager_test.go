package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-tracker/internal/api/handlers"
	"task-tracker/internal/models"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateTask_AsManager_ShouldReturnForbidden(t *testing.T) {
	db := setupTestDB(t)

	taskInput := handlers.TaskRequest{
		Summary: "Manager task",
		Date:    time.Now(),
	}
	body, _ := json.Marshal(taskInput)

	user := getManagerUser()
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.CreateTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestUpdateTask_AsManager_ShouldReturnForbidden(t *testing.T) {
	db := setupTestDB(t)

	task := models.Task{Summary: "Old summary", UserID: 1, Date: time.Now()}
	db.Create(&task)

	updated := handlers.TaskRequest{
		Summary: "Updated summary",
		Date:    time.Now().AddDate(0, 0, 1),
	}
	body, _ := json.Marshal(updated)

	user := getManagerUser()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(task.ID)})
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.UpdateTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestListTasks_AsManager_ShouldReturnOK(t *testing.T) {
	db := setupTestDB(t)

	tasks := []models.Task{
		{Summary: "Task 1", UserID: 1, Date: time.Now()},
		{Summary: "Task 2", UserID: 2, Date: time.Now()},
	}
	for _, task := range tasks {
		db.Create(&task)
	}

	user := getManagerUser()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.ListTasks(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got []models.Task
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestDeleteTask_AsManager_ShouldReturnOK(t *testing.T) {
	db := setupTestDB(t)

	task := models.Task{Summary: "To be deleted", UserID: 1, Date: time.Now()}
	db.Create(&task)

	user := getManagerUser()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(task.ID)})
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.DeleteTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}
