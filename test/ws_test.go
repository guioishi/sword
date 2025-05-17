package test

import (
	"net/http"
	"net/http/httptest"
	"task-tracker/internal/api/handlers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWS_AsTechnician_ShouldReturnForbidden(t *testing.T) {
	user := getTechUser()
	req := httptest.NewRequest(http.MethodGet, "/ws/notifications", nil)
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.NotificationWebSocket()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// Return BadRequest because of missing of required headers for websockets
func TestWS_AsManager_ShouldReturnBadRequest(t *testing.T) {
	user := getManagerUser()
	req := httptest.NewRequest(http.MethodGet, "/ws/notifications", nil)
	req = mockContextWithUser(req, user)

	w := httptest.NewRecorder()
	handler := handlers.NotificationWebSocket()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
