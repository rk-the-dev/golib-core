package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

// KafkaClient defines the interface for Kafka producer & consumer
type KafkaClient interface {
	Produce(ctx context.Context, topic string, key, message []byte) error
	Consume(ctx context.Context, topic string, groupID string, handler func(message kafka.Message)) error
	Close() error
}

// KafkaConfig defines Kafka connection configurations
type KafkaConfig struct {
	Brokers   []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ClientID  string   `env:"KAFKA_CLIENT_ID" envDefault:"golang-client"`
	Partition int      `env:"KAFKA_PARTITION" envDefault:"0"`
}

// kafkaClient implements KafkaClient
type kafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

var (
	instance *kafkaClient
	once     sync.Once
)

// NewKafkaClient initializes and returns a Kafka producer
func NewKafkaClient(cfg *KafkaConfig) KafkaClient {
	once.Do(func() {
		instance = &kafkaClient{
			writer: &kafka.Writer{
				Addr:     kafka.TCP(cfg.Brokers...),
				Balancer: &kafka.LeastBytes{},
			},
		}
		fmt.Println("✅ Kafka client initialized:", cfg.Brokers)
	})

	return instance
}

// Produce sends a message to Kafka
func (k *kafkaClient) Produce(ctx context.Context, topic string, key, message []byte) error {
	err := k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
	})
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}
	fmt.Println("✅ Message sent to Kafka:", string(message))
	return nil
}

// Consume reads messages from Kafka and processes them using the provided handler
func (k *kafkaClient) Consume(ctx context.Context, topic, groupID string, handler func(message kafka.Message)) error {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	fmt.Println("🔄 Consuming messages from topic:", topic)

	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("error reading message from Kafka: %w", err)
		}
		handler(msg) // Process the message using the custom handler
	}
}

// Close closes Kafka producer & consumer connections
func (k *kafkaClient) Close() error {
	if k.writer != nil {
		k.writer.Close()
	}
	if k.reader != nil {
		k.reader.Close()
	}
	fmt.Println("🔻 Kafka client closed")
	return nil
}
