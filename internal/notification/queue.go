package notification

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var queue amqp.Queue
var ch *amqp.Channel

func InitQueue() error {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	queue, err = ch.QueueDeclare("task_notifications", true, false, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}

type TaskNotification struct {
	TechID  uint   `json:"tech_id"`
	Summary string `json:"summary"`
	Date    string `json:"date"`
}

func PublishNotification(n TaskNotification) {
	if ch == nil {
		return
	}
	body, _ := json.Marshal(n)
	ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func ConsumeStream() (<-chan amqp.Delivery, error) {
	if ch == nil {
		return nil, fmt.Errorf("RabbitMQ channel not initialized")
	}
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
