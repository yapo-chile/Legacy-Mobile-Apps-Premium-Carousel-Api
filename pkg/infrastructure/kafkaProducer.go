package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/Yapo/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka" // nolint
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/interfaces/repository"
)

// KafkaProducer struct representing a message producer for kafka
type KafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new KafkaProducer with the given brokers
func NewKafkaProducer(
	host string,
	port int,
	acks string,
	compressionType string,
	retries int,
	lingerMS int,
	requestTimeoutMS int,
	enableIdempotence bool,

) (repository.KafkaProducer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%v:%d", host, port),
		"acks":               acks,
		"compression.type":   compressionType,
		"retries":            retries,
		"linger.ms":          lingerMS,
		"request.timeout.ms": requestTimeoutMS,
		"enable.idempotence": enableIdempotence,
	}
	producer, err := kafka.NewProducer(conf)
	if err != nil {
		logger.Crit("Failed to create producer: %s\n", err)
		return nil, err
	}
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		logger.Info("Producer connected to kafka using config: \n%s\n", jconf)
	} else {
		logger.Info("Producer connected to kafka using config: \n%+v\n", conf)
	}
	return &KafkaProducer{producer: producer}, nil
}

// SendMessage sends a message with the specified topic
func (k KafkaProducer) SendMessage(topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)
	e := <-deliveryChan
	m := e.(*kafka.Message)
	err = m.TopicPartition.Error
	if err != nil {
		logger.Error("Failed to send the message %s: %v", string(message), err)
	} else {
		logger.Info("Delivered message to topic %s [%d] at offset %v",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
	return err
}

// Close close the KafkaProducer
func (k KafkaProducer) Close() error {
	k.producer.Close()
	return nil
}
