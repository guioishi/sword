package main

import (
	"log"
	"net/http"
	"os"
	"task-tracker/internal/api"
	"task-tracker/internal/db"
	"task-tracker/internal/notification"
)

func main() {
	dbConn := db.Init()
	router := api.SetupRouter(dbConn)

	err := notification.InitQueue()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	http.ListenAndServe(":"+port, router)
}
