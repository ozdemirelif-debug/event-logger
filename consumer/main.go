package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type Event struct {
	EventID    string                 `json:"eventId"`
	Source     string                 `json:"source"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	Timestamp  string                 `json:"timestamp"`
	ReceivedAt time.Time              `json:"receivedAt"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	connStr := "postgres://eventuser:eventpass@postgres:5432/eventsdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	failOnError(err, "PostgreSQL'e bağlanılamadı")
	defer db.Close()

	err = db.Ping()
	failOnError(err, "Veritabanı ping başarısız")

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS events (
		event_id TEXT PRIMARY KEY,
		source TEXT,
		type TEXT,
		payload JSONB,
		timestamp TIMESTAMP,
		received_at TIMESTAMP	
	);
	`)
	failOnError(err, "Tablo oluşturulamadı")

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "RabbitMQ'ya bağlanılamadı")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Kanal açılamadı")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"event-queue", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Kuyruk oluşturulamadı")
	msgs, err := ch.Consume(
		"event-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Mesaj dinleme başlatılamadı")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var event Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Geçersiz event formatı: %s", err)
				continue
			}

			event.ReceivedAt = time.Now()

			eventTime, err := time.Parse(time.RFC3339, event.Timestamp)
			if err != nil {
				log.Printf("Timestamp parse edilemedi: %s", err)
				continue
			}
			payloadJSON, err := json.Marshal(event.Payload)
			if err != nil {
				log.Printf("Payload JSON'a çevrilemedi: %s", err)
				continue
			}

			_, err = db.Exec(
				`INSERT INTO events(event_id, source, type, payload, timestamp, received_at) VALUES ($1,$2,$3,$4,$5,$6)`,
				event.EventID,
				event.Source,
				event.Type,
				payloadJSON,
				eventTime,
				event.ReceivedAt,
			)
			if err != nil {
				log.Printf("Veritabanına kaydedilemedi: %s", err)
				continue
			}

			log.Printf("Event kaydedildi: %s", event.EventID)
		}
	}()

	log.Println("Consumer dinleniyor...")
	<-forever
}
