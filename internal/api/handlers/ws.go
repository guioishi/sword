package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-tracker/internal/auth"
	"task-tracker/internal/notification"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NotificationWebSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		user, _ := auth.GetUser(r.Context())
		if user.Role != "manager" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		msgs, err := notification.ConsumeStream()
		if err != nil {
			fmt.Println("Failed to consume RabbitMQ:", err)
			return
		}

		for msg := range msgs {
			var notif notification.TaskNotification
			if err := json.Unmarshal(msg.Body, &notif); err != nil {
				continue
			}
			err = conn.WriteJSON(notif)
			if err != nil {
				fmt.Println("Write to WebSocket failed:", err)
				break
			}
		}
	}
}
