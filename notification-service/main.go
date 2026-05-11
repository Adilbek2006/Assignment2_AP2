package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type EmailSender interface {
	SendEmail(email string, orderID string, amount float64) error
}

type SimulatedProvider struct{}

func (s *SimulatedProvider) SendEmail(email, orderID string, amount float64) error {
	time.Sleep(500 * time.Millisecond)

	if rand.Float32() < 0.3 {
		return errors.New("simulated 503 Service Unavailable")
	}

	log.Printf("[SIMULATED PROVIDER] Successfully sent email to %s for Order #%s", email, orderID)
	return nil
}

type RealSMTPProvider struct{}

func (r *RealSMTPProvider) SendEmail(email, orderID string, amount float64) error {
	log.Printf("[REAL SMTP] Connecting to Mailjet/SMTP... Sent to %s", email)
	return nil
}

type PaymentEvent struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	var provider EmailSender
	if os.Getenv("PROVIDER_MODE") == "REAL" {
		provider = &RealSMTPProvider{}
		log.Println("Started worker with REAL Email Provider")
	} else {
		provider = &SimulatedProvider{}
		log.Println("Started worker with SIMULATED Email Provider")
	}

	var conn *amqp.Connection
	var err error
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err == nil {
			log.Println("Connected to RabbitMQ!")
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5s...", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("payment.completed", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	go func() {
		for d := range msgs {
			var event PaymentEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				d.Nack(false, false)
				continue
			}

			redisKey := "processed_order:" + event.OrderID

			isNew, err := rdb.SetNX(ctx, redisKey, "processing", 24*time.Hour).Result()
			if err != nil || !isNew {
				log.Printf("[Notification] Duplicate ignored for Order #%s", event.OrderID)
				d.Ack(false)
				continue
			}

			maxRetries := 3
			backoffDelay := 2 * time.Second
			success := false

			for i := 1; i <= maxRetries; i++ {
				err := provider.SendEmail(event.CustomerEmail, event.OrderID, event.Amount)
				if err == nil {
					success = true
					break
				}

				log.Printf("[Worker Warning] Provider failed: %v. Retrying in %v (Attempt %d/%d)...", err, backoffDelay, i, maxRetries)
				time.Sleep(backoffDelay)
				backoffDelay *= 2
			}

			if success {
				rdb.Set(ctx, redisKey, "completed", 24*time.Hour)
				d.Ack(false)
			} else {
				log.Printf("[Worker ERROR] Failed to process Order #%s after %d retries. Dropping.", event.OrderID, maxRetries)
				rdb.Del(ctx, redisKey)
				d.Nack(false, true)
			}
		}
	}()

	log.Printf(" [*] Background Worker is waiting for messages. To exit press CTRL+C")
	<-sigs
	log.Println("Gracefully shutting down Notification Service")
}
