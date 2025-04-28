package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received termination signal. Closing consumer...")
		cancel()
	}()

	// Set up Kafka reader
	reader := getKafkaReader()
	defer reader.Close()

	log.Println("Consumer started, waiting for messages...")

	// Read messages in a loop
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					// Context was canceled, exit gracefully
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			log.Printf("Message received from Kafka: %s ::Time %s ", string(message.Value), message.Time)
			

		}
	}
}

func getKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{"localhost:9092"},
		Topic:            "order-topic",
		GroupID:          "khach-hang-vip",
		MinBytes:         1e3, 
		MaxBytes:         1e6,   
		MaxWait:          100 * time.Millisecond, 
		ReadBatchTimeout: 100 * time.Millisecond,
		CommitInterval:   0, 
		StartOffset:      kafka.FirstOffset,
		Partition:        0,
	})
}