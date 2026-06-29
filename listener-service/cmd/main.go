package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Listener Service")

	// #1 try to connect to RabbitMQ (with re-try)
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// #2 start listening to messages
	log.Println("Listening and consuming for RabbitMQ messages ...")

	// create consumer

	// watch the queue amd consume events
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err == nil {
			log.Println("Connected to RabbitMQ")
			return conn, nil
		}

		log.Printf("RabbitMQ not ready yet (attempt %d/%d): %v", i, maxRetries, err)

		backOff := time.Duration(i*i) * time.Second
		log.Printf("Backing off for %v...", backOff)

		time.Sleep(backOff)
	}

	return nil, err
}
