package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"strconv"
	"time"
)

func CreateAdminClient(brokers string) *kafka.AdminClient {

	client, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil
	}
	return client
}

func (k *Kafka) GetKafkaMetadata() (*kafka.Metadata, error) {
	res, err := k.AdminClient.GetMetadata(nil, true, 60)
	if err != nil {
		k.Log.Errorf("Failed to get metadata information from kafka cluster: %s\n", err)
		return nil, err
	}
	return res, nil
}

func (k *Kafka) GetKafkaTopicMetadata(topic string) (*kafka.Metadata, error) {
	res, err := k.AdminClient.GetMetadata(&topic, true, 60)
	if err != nil {
		k.Log.Errorf("Failed to get metadata information about the topic from kafka cluster: %s\n", err)
		return nil, err
	}
	return res, nil
}

func (k *Kafka) CreateTopic(topic string, numParts string, replicationFactor string) ([]kafka.TopicResult, error) {
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		return nil, err
	}
	parititons, err := strconv.Atoi(numParts)
	if err != nil {
		k.Log.Errorf("Invalid partition count: %s: %v\n", parititons, err)
		return nil, err
	}
	rFactor, err := strconv.Atoi(replicationFactor)
	if err != nil {
		k.Log.Errorf("Invalid replication factor: %s: %v\n", replicationFactor, err)
		return nil, err
	}
	results, err := k.AdminClient.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     parititons,
			ReplicationFactor: rFactor}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))

	if err != nil {
		k.Log("Failed to create topic: %v\n", err)
		return nil, err
	}

	return results, nil
}

func (k *Kafka) DeleteKafkaTopic(topic string) ([]kafka.TopicResult, error) {
	topicSplice := []string{topic}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete topics on cluster
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		return nil, err
	}

	results, err := k.AdminClient.DeleteTopics(ctx, topicSplice, kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		k.Log.Errorf("Failed to delete topics: %v\n", err)
		return nil, err
	}

	return results, nil
}

func (k *Kafka) DescribeConfig(resourceType, resourceName string) ([]kafka.ConfigResourceResult, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaResourceType, err := kafka.ResourceTypeFromString(resourceType)
	if err != nil {
		k.Log.Errorf("Invalid resource type: %s\n", os.Args[2])
		return nil, err
	}

	dur, _ := time.ParseDuration("20s")

	results, err := k.AdminClient.DescribeConfigs(ctx,
		[]kafka.ConfigResource{{Type: kafkaResourceType, Name: resourceName}},
		kafka.SetAdminRequestTimeout(dur))
	if err != nil {
		k.Log.Errorf("Failed to DescribeConfigs(%s, %s): %s\n",
			resourceType, resourceName, err)
		return nil, err
	}

	return results, nil
}
