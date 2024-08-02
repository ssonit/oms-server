package kafka

import (
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

const (
	OrderCreatedEvent = "order.created"
	OrderPaidEvent    = "order.paid"
)

var (
	kafkaURL = "localhost:9092"
	topic    = "orders"
)

// producer
func GetKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// consumer
func GetKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})
}
