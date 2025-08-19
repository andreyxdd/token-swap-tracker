package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"producer/internal/simulator"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaConfig holds Kafka producer configuration
type KafkaConfig struct {
	Brokers   []string
	Topic     string
	Partition int32
}

var swapChannel = make(chan *simulator.SwapEvent, 1000) // Swap event channel (buffered)
var kafkaConfig = KafkaConfig{
	Brokers:   []string{"kafka:9093"}, // Kafka broker address
	Topic:     "swaps",                // Kafka topic for swap events
	Partition: kafka.PartitionAny,     // Any partition (you can customize this)
}

func main() {
	// Graceful shutdown signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Initialize Kafka Producer
	producer, err := initKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()

	// Start the Kafka producer loop in a goroutine
	go produceToKafka(producer)

	// Simulate receiving swap events and pushing to the channel
	go simulator.SimulateSwapEvents(swapChannel)

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down gracefully...")
}

// Initialize Kafka producer
func initKafkaProducer() (*kafka.Producer, error) {
	brokers := strings.Join(kafkaConfig.Brokers, ",")
	config := kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "all",  // Ensure that all replicas acknowledge
		"compression.codec": "gzip", // Optional: compress messages to save bandwidth
		"linger.ms":         5,      // Delay for batching to increase throughput
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	// Log any errors from Kafka producer
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Error producing message: %v", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return producer, nil
}

// Function to push events from the channel to Kafka
func produceToKafka(producer *kafka.Producer) {
	for event := range swapChannel {
		// Marshal the SwapEvent to JSON
		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error marshalling event: %v", err)
			continue
		}

		// Create a Kafka message
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &kafkaConfig.Topic, Partition: kafkaConfig.Partition},
			Value:          eventBytes,
		}

		// Produce message to Kafka
		err = producer.Produce(msg, nil)
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
}
