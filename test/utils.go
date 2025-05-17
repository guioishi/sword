package test

import (
	"context"
	"net/http"
	"task-tracker/internal/auth"
	"task-tracker/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Task{}, &models.User{})
	require.NoError(t, err)
	return db
}

func mockContextWithUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), auth.UserContextKey, auth.UserContext{ID: user.ID, Role: user.Role})
	return r.WithContext(ctx)
}

func getManagerUser() *models.User {
	return &models.User{ID: 99, Role: "manager"}
}

func getTechUser() *models.User {
	return &models.User{ID: 1, Role: "technician"}
}
