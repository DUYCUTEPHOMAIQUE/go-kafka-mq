package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	writer := getKafkaWriter()
	defer writer.Close()

	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		//tao msg
		msgKK := kafka.Message{
			Value: []byte(msg.Content),
			Time: time.Now(),
		}
		// send msg
		err := writer.WriteMessages(context.Background(), msgKK)
		if err != nil {
			log.Printf("Failed to write message: %v", err)
			http.Error(w, "Failed to process message", http.StatusInternalServerError)
			return
		}

		log.Printf("Message sent to Kafka: %s", msg.Content)
	})

	log.Println("Producer server starting on :8999...")
	if err := http.ListenAndServe(":8999", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order-topic",
		Balancer: &kafka.LeastBytes{},
	}
}