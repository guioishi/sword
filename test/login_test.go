package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-tracker/internal/api/handlers"
	"task-tracker/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func registerUsers(db *gorm.DB) {
	password := "1234"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	db.Create(&models.User{Username: "tech", Password: string(hashedPass), Role: "technician"})
	db.Create(&models.User{Username: "manager", Password: string(hashedPass), Role: "manager"})
}

func TestLogin_AsTechnician_ShouldHaveRoleTechnician(t *testing.T) {
	db := setupTestDB(t)
	registerUsers(db)

	bodyInput := handlers.LoginRequest{
		Username: "tech",
		Password: "1234",
	}
	body, _ := json.Marshal(bodyInput)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler := handlers.Login(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, "technician", result["role"])
}

func TestLogin_AsManager_ShouldHaveRoleManager(t *testing.T) {
	db := setupTestDB(t)
	registerUsers(db)

	bodyInput := handlers.LoginRequest{
		Username: "manager",
		Password: "1234",
	}
	body, _ := json.Marshal(bodyInput)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler := handlers.Login(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, "manager", result["role"])
}

func TestLogin_WithWrongUser_ShouldReturnUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	registerUsers(db)

	bodyInput := handlers.LoginRequest{
		Username: "wrong",
		Password: "1234",
	}
	body, _ := json.Marshal(bodyInput)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler := handlers.Login(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLogin_WithWrongPass_ShouldReturnUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	registerUsers(db)

	bodyInput := handlers.LoginRequest{
		Username: "manager",
		Password: "wrong",
	}
	body, _ := json.Marshal(bodyInput)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler := handlers.Login(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
