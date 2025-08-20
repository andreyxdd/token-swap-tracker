package main

import (
	"consumer/internal/consumer"
	"consumer/internal/services"
	"consumer/internal/ws"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	wsPort := os.Getenv("WS_PORT")
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	redisCfg := services.RedisConfig{Addr: redisAddr, Password: redisPw}
	repo := services.NewRedisStatsRepo(redisCfg)
	service := services.NewStatsService(repo)

	var wsCh = make(chan []byte)
	c, err := consumer.New(service, cfg, wsCh, sigCh)
	if err != nil {
		log.Printf("failed to initialize Kafka consumer: %v\n", err)
	}
	defer c.KafkaConsumer.Close()

	ws := ws.New(wsCh)
	http.HandleFunc("/ws", ws.Handler)
	addr := fmt.Sprintf(":%s", wsPort)

	go func() {
		log.Printf("WebSocket server started on %s\n", addr)
		err = http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("Error starting WebSocket server: %v\n", err)
		}
		log.Println("WebSocket server stopped")
	}()

	// Start web-socket broadcasting and consuming events from kafka
	go ws.HandleBroadcasting()
	go c.ProcessSwapEvents()

	<-sigCh
	log.Println("Shutting down gracefully...")
}
