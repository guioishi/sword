package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"task-tracker/internal/auth"
	"task-tracker/internal/models"
	"task-tracker/internal/notification"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TaskRequest struct {
	Summary string    `json:"summary"`
	Date    time.Time `json:"date"`
}

func CreateTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetUser(r.Context())
		if user.Role != "technician" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var input TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		task := models.Task{
			Summary: input.Summary,
			Date:    input.Date,
			UserID:  user.ID,
		}
		db.Create(&task)
		go func() {
			if user.Role == "technician" {
				notification.PublishNotification(notification.TaskNotification{
					TechID:  task.UserID,
					Summary: task.Summary,
					Date:    task.Date.Format("2006-01-02"),
				})
			}
		}()
		json.NewEncoder(w).Encode(task)
	}
}

func ListTasks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetUser(r.Context())
		var tasks []models.Task
		if user.Role == "manager" {
			db.Table("tasks").Find(&tasks)
		} else {
			db.Where("user_id = ?", user.ID).Find(&tasks)
		}
		json.NewEncoder(w).Encode(tasks)
	}
}

func UpdateTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetUser(r.Context())
		if user.Role != "technician" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		id := mux.Vars(r)["id"]
		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if task.UserID != user.ID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var input TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		db.Model(&task).Updates(models.Task{Summary: input.Summary, Date: input.Date})

		go func() {
			if user.Role == "technician" {
				notification.PublishNotification(notification.TaskNotification{
					TechID:  task.UserID,
					Summary: task.Summary,
					Date:    task.Date.Format("2006-01-02"),
				})
			}
		}()

		json.NewEncoder(w).Encode(task)
	}
}

func DeleteTask(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetUser(r.Context())
		if user.Role != "manager" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		id := mux.Vars(r)["id"]
		db.Delete(&models.Task{}, id)
		w.WriteHeader(http.StatusNoContent)
	}
}
