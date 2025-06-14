package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

type Event struct {
	Source    string                 `json:"source"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "RabbitMQ'ya bağlanılamadı")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Kanal açılamadı")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"event-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Kuyruk oluşturulamadı")

	router := gin.Default()

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Basic Auth middleware
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password123",
	}))

	authorized.POST("/events", func(c *gin.Context) {
		var event Event
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz JSON"})
			return
		}

		eventID := uuid.New().String()
		fullEvent := map[string]interface{}{
			"eventId":   eventID,
			"source":    event.Source,
			"type":      event.Type,
			"payload":   event.Payload,
			"timestamp": event.Timestamp,
		}

		body, err := json.Marshal(fullEvent)
		failOnError(err, "JSON'a dönüştürülemedi")

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
		failOnError(err, "Mesaj gönderilemedi")

		c.JSON(http.StatusOK, gin.H{
			"status":  "queued",
			"eventId": eventID,
		})
	})

	router.Run(":8080")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}