package main

import (
	"log"
	"os"
	"os/signal"
	"producer/internal/producer"
	"producer/internal/simulator"
	"strconv"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	swapEventsPerSecondStr := os.Getenv("SWAP_EVENTS_PER_SECOND")
	swapEventsPerSecond, err := strconv.ParseFloat(swapEventsPerSecondStr, 64)
	if err != nil {
		log.Fatalf("Error converting string to float64: %v", err)
	}

	var swapChannel = make(chan *simulator.SwapEvent, 1000)
	var cfg = producer.Config{
		Brokers: kafkaBrokers,
		Topic:   kafkaTopic,
	}

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	p, err := producer.New(swapChannel, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer p.KafkaProducer.Close()

	// Log any errors from Kafka producer
	go p.LogErrors()
	// and start the Kafka producer loop in a goroutine
	go p.ProduceMessagesFromChannel()

	// Simulate receiving swap events by pushing to the channel
	sim := simulator.New(swapChannel, swapEventsPerSecond)
	go sim.SimulateSwapEvents()

	<-sigCh
	log.Println("Shutting down gracefully...")
}
