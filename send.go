package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Event struct {
	EventID   string                 `json:"eventId"`
	Source    string                 `json:"source"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "RabbitMQ bağlantısı kurulamadı")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Kanal oluşturulamadı")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"event-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Queue oluşturulamadı")

	event := Event{
		EventID:   "event-123",
		Source:    "api",
		Type:      "test-event",
		Payload:   map[string]interface{}{"message": "Merhaba RabbitMQ!"},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(event)
	failOnError(err, "Event JSON'a dönüştürülemedi")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	failOnError(err, "Mesaj yayınlanamadı")

	log.Println("Event gönderildi:", string(body))
}
