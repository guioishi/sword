package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-tracker/internal/api/handlers"
	"task-tracker/internal/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateTask_AsTechnician_ShouldReturnOK(t *testing.T) {
	db := setupTestDB(t)

	taskInput := handlers.TaskRequest{
		Summary: "Test task",
		Date:    time.Now(),
	}
	body, _ := json.Marshal(taskInput)

	user := getTechUser()
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.CreateTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var task models.Task
	err := json.NewDecoder(resp.Body).Decode(&task)
	require.NoError(t, err)
	require.Equal(t, taskInput.Summary, task.Summary)
	require.Equal(t, user.ID, task.UserID)
}

func TestUpdateOwnTask_AsTechnician_ShouldReturnOK(t *testing.T) {
	db := setupTestDB(t)

	user := getTechUser()
	task := models.Task{Summary: "Old summary", UserID: user.ID, Date: time.Now()}
	db.Create(&task)

	updated := handlers.TaskRequest{
		Summary: "Updated summary",
		Date:    time.Now().AddDate(0, 0, 1),
	}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(task.ID)})
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.UpdateTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.Task
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, updated.Summary, result.Summary)
}

func TestUpdateAnotherTask_AsTechnician_ShouldReturnForbidden(t *testing.T) {
	db := setupTestDB(t)

	user := getTechUser()
	task := models.Task{Summary: "Old summary", UserID: user.ID + 1, Date: time.Now()}
	db.Create(&task)

	updated := handlers.TaskRequest{
		Summary: "Updated summary",
		Date:    time.Now().AddDate(0, 0, 1),
	}
	body, _ := json.Marshal(updated)

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

func TestListOwnTasks_AsTechnician_ShouldReturnOK(t *testing.T) {
	db := setupTestDB(t)

	user := getTechUser()
	tasks := []models.Task{
		{Summary: "Task 1", UserID: user.ID, Date: time.Now()},
		{Summary: "Task 2", UserID: user.ID + 1, Date: time.Now()},
	}
	for _, task := range tasks {
		db.Create(&task)
	}

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
	require.Len(t, got, 1)
	require.Equal(t, user.ID, got[0].UserID)
}

func TestDeleteTask_AsTechnician_ShouldReturnForbidden(t *testing.T) {
	db := setupTestDB(t)

	user := getTechUser()
	task := models.Task{Summary: "To be deleted", UserID: user.ID, Date: time.Now()}
	db.Create(&task)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(task.ID)})
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.DeleteTask(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}
