package api

import (
	corsHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"task-tracker/internal/api/handlers"
	"task-tracker/internal/auth"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	corsHandler := corsHandlers.CORS(
		corsHandlers.AllowedOrigins([]string{"*"}),
		corsHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		corsHandlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
	)

	loginRouter := r.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("", handlers.Login(db)).Methods("POST", "OPTIONS")

	taskRouter := r.PathPrefix("/tasks").Subrouter()
	taskRouter.Use(auth.Middleware)
	taskRouter.HandleFunc("", handlers.CreateTask(db)).Methods("POST")
	taskRouter.HandleFunc("", handlers.ListTasks(db)).Methods("GET")
	taskRouter.HandleFunc("/{id}", handlers.UpdateTask(db)).Methods("PUT")
	taskRouter.HandleFunc("/{id}", handlers.DeleteTask(db)).Methods("DELETE")

	wsRouter := r.PathPrefix("/ws").Subrouter()
	wsRouter.Use(auth.Middleware)
	wsRouter.HandleFunc("/notifications", handlers.NotificationWebSocket()).Methods("GET")

	r.Use(corsHandler)

	return r
}
