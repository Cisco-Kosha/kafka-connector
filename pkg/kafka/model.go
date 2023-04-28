package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kosha/kafka-connector/pkg/config"
	logger "github.com/kosha/kafka-connector/pkg/logger"
)

type Kafka struct {
	Brokers     string             `json:"brokers,omitempty"`
	AdminClient *kafka.AdminClient `json:"adminClient,omitempty"`
	Log         logger.Logger      `json:"log,omitempty"`
	Producer    *kafka.Producer    `json:"producer,omitempty"`
	Consumers   []*Consumer        `json:"consumer,omitempty"`
}

type Topic struct {
	Topic             string `json:"topic,omitempty"`
	NumParts          string `json:"numParts,omitempty"`
	ReplicationFactor string `json:"replicationFactor,omitempty"`
}

type MetadataTopicJson struct {
	Topic       string `json:"topic"`
	PartitionID string `json:"partitionId"`
}

func NewKafkaClient(cfg *config.Config, log logger.Logger) *Kafka {

	k := Kafka{}

	k.Brokers = cfg.GetBrokerAddress()
	k.Log = log
	k.AdminClient = CreateAdminClient(k.Brokers)

	if k.AdminClient == nil {
		panic("Could not create Admin Client")
	} else {
		k.Log.Info("Successfully created an Kafka Admin Client")
	}

	k.Producer = CreateProducer(k.Brokers)

	k.Consumers = []*Consumer{}

	return &k
}

type Client interface {
	// GetKafkaMetadata Get information about all kafka topics
	GetKafkaMetadata() *kafka.Metadata
}
