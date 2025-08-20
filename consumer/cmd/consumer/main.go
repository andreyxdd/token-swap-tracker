package main

import (
	"consumer/internal/consumer"
	"consumer/internal/services"
	"log"
	"os"
	"os/signal"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaConsumerGroupID := os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPw := os.Getenv("REDIS_PASSWORD")
	debugStr := os.Getenv("DEBUG")
	debug := false
	if debugStr == "true" {
		debug = true
	}

	var cfg = consumer.Config{
		Brokers: kafkaBrokers,
		Topic:   kafkaTopic,
		GroupId: kafkaConsumerGroupID,
		Debug:   debug,
	}

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	redisCfg := services.RedisConfig{Addr: redisAddr, Password: redisPw}
	repo := services.NewRedisStatsRepo(redisCfg)
	service := services.NewStatsService(repo)
	c, err := consumer.New(service, cfg)
	if err != nil {
		log.Fatalf("failed to initialize Kafka consumer: %v", err)
	}
	defer c.KafkaConsumer.Close()

	// Start consuming events
	go c.ProcessSwapEvents(sigCh)

	<-sigCh
	log.Println("Shutting down gracefully...")
}
