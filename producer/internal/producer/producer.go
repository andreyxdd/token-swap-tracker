package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client[T any] struct {
	ch            chan T
	cfg           Config
	KafkaProducer *kafka.Producer
}

// Initialize Kafka producer
func New[T any](ch chan T, cfg Config) (*Client[T], error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"acks":              "all",  // Ensure that all replicas acknowledge
		"compression.codec": "gzip", // Optional: compress messages to save bandwidth
		"linger.ms":         5,      // Delay for batching to increase throughput
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return &Client[T]{ch: ch, cfg: cfg, KafkaProducer: producer}, nil
}

func (c *Client[T]) ProduceMessagesFromChannel() {
	for event := range c.ch {
		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal an event: %v", err)
			continue
		}

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &c.cfg.Topic},
			Value:          eventBytes,
		}

		err = c.KafkaProducer.Produce(msg, nil)
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}
}

func (c *Client[T]) LogErrors() {
	for e := range c.KafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Error producing message: %v", ev.TopicPartition.Error)
			}
		}
	}
}
