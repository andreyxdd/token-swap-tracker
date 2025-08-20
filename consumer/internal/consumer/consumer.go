package consumer

import (
	"consumer/internal/models"
	"consumer/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
)

type Client struct {
	statsService  *services.StatsService
	cfg           Config
	KafkaConsumer *kafka.Consumer
}

func New(statsService *services.StatsService, cfg Config) (*Client, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka consumer")
	}

	err = consumer.Subscribe(cfg.Topic, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to a topic %s", cfg.Topic)
	}
	return &Client{statsService, cfg, consumer}, nil
}

// Consume from Kafka and process stats
func (c *Client) ProcessSwapEvents(sigCh chan os.Signal) {
	for {
		select {
		case <-sigCh:
			log.Println("Received shutdown signal, stopping consumer...")
			return
		default:
			msg, err := c.KafkaConsumer.ReadMessage(-1)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				log.Fatal("failed to read message with kafka consumer: ", err)
				continue
			}

			var event models.SwapEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Fatal("failed to unmarshal message: ", err)
			}

			err = c.statsService.ProcessSwapEvent(context.Background(), event)
			if err != nil {
				log.Fatalf("failed to process swap event with tx hash %s: %v", event.TxHash, err)
			}

			if c.cfg.Debug {
				logSwapEvent(event, msg)
			}
		}
	}
}

func logSwapEvent(event models.SwapEvent, msg *kafka.Message) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 90))
	fmt.Printf("📊 SWAP EVENT TX HASH %s\n", event.TxHash)
	fmt.Printf("%s\n", strings.Repeat("=", 90))

	fmt.Printf("🔧 Kafka Metadata:\n")
	fmt.Printf("   Topic: %s | Partition: %d | Offset: %d\n",
		*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	fmt.Printf("   Timestamp: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"))

	fmt.Printf("\n💱 Swap Details:\n")
	fmt.Printf("   %s → %s\n", event.TokenFrom, event.TokenTo)
	fmt.Printf("   Amount: %.6f %s → %.6f %s\n",
		event.AmountFrom, event.TokenFrom, event.AmountTo, event.TokenTo)

	fmt.Printf("\n💰 USD Values:\n")
	fmt.Printf("   $%.2f\n", event.UsdValue)

	fmt.Printf("\n📈 Trading Info:\n")
	fmt.Printf("   Event Time: %s\n", event.Timestamp.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("   Processing Delay: %v\n", time.Since(event.Timestamp).Truncate(time.Millisecond))

	fmt.Printf("%s\n", strings.Repeat("=", 90))
}
