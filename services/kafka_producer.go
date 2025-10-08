package services

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hash"
	"strings"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/xdg-go/scram"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

var kafkaProducerInstance *KafkaProducer

// SCRAM client generator for SCRAM-SHA-256
var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

// InitKafkaProducer initializes the Kafka producer with configurable security
func InitKafkaProducer() error {
	brokers := strings.Split(viper.GetString("KAFKA_BOOTSTRAP_SERVERS"), ",")
	topic := viper.GetString("KAFKA_TOPIC_EXECUTION_REQUESTS")
	securityProtocol := viper.GetString("KAFKA_SECURITY_PROTOCOL")

	// Default to SASL_SSL if not specified
	if securityProtocol == "" {
		securityProtocol = "SASL_SSL"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// Configure security based on protocol
	switch strings.ToUpper(securityProtocol) {
	case "SASL_SSL":
		// SASL_SSL configuration for Aiven Kafka
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.User = viper.GetString("KAFKA_SASL_USERNAME")
		config.Net.SASL.Password = viper.GetString("KAFKA_SASL_PASSWORD")
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification for Aiven cloud
		}
		log.Info().Msg("Kafka producer configured with SASL_SSL security")

	case "PLAINTEXT":
		// PLAINTEXT configuration for local Kafka
		config.Net.SASL.Enable = false
		config.Net.TLS.Enable = false
		log.Info().Msg("Kafka producer configured with PLAINTEXT security")

	default:
		log.Error().Str("protocol", securityProtocol).Msg("Unsupported Kafka security protocol")
		return fmt.Errorf("unsupported Kafka security protocol: %s", securityProtocol)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Kafka producer")
		return err
	}

	kafkaProducerInstance = &KafkaProducer{
		producer: producer,
		topic:    topic,
	}

	log.Info().Str("topic", topic).Msg("Kafka producer initialized successfully")
	return nil
}

// GetKafkaProducer returns the singleton Kafka producer instance
func GetKafkaProducer() *KafkaProducer {
	return kafkaProducerInstance
}

// ExecutionRequest represents the message payload for report execution
type ExecutionRequest struct {
	ExecutionID string `json:"execution_id"`
	ConfigID    int    `json:"config_id"`
	ScheduleID  *int   `json:"schedule_id"`
	ExecutedBy  string `json:"executed_by"`
	QueuedAt    string `json:"queued_at"`
}

// ProduceExecutionRequest sends a report execution request to Kafka
func (kp *KafkaProducer) ProduceExecutionRequest(req ExecutionRequest) error {
	messageJSON, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal execution request")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Key:   sarama.StringEncoder(req.ExecutionID),
		Value: sarama.ByteEncoder(messageJSON),
	}

	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		log.Error().
			Err(err).
			Str("execution_id", req.ExecutionID).
			Msg("Failed to produce message to Kafka")
		return err
	}

	log.Info().
		Str("execution_id", req.ExecutionID).
		Int32("partition", partition).
		Int64("offset", offset).
		Msg("Successfully produced execution request to Kafka")

	return nil
}

// Close closes the Kafka producer connection
func (kp *KafkaProducer) Close() error {
	if kp.producer != nil {
		return kp.producer.Close()
	}
	return nil
}
