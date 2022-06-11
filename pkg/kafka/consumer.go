package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Name     string
	Offset   string
	Consumer *kafka.Consumer
}

func (k *Kafka) CreateConsumer(name, offset string) error {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.Brokers,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              name,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     offset})

	if err != nil {
		k.Log.Errorf("Failed to create consumer: %s\n", err)
		return nil
	}

	k.Consumers = append(k.Consumers, &Consumer{
		Name:     name,
		Offset:   offset,
		Consumer: c,
	})
	k.Log.Infof("Created Consumer %v\n", c)
	return nil
}

func (k *Kafka) SubscribeToTopics(topic string, kafkaConsumerName string) error {

	consumers := k.Consumers

	var c *kafka.Consumer

	for _, consumer := range consumers {
		if kafkaConsumerName == consumer.Name {
			c = consumer.Consumer
		}
	}

	err := c.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	k.Log.Infof("Subscribed to topic %v\n", topic)
	return nil
}

func (k *Kafka) ConsumeMessage(kafkaConsumerName string) (string, string, error) {

	consumers := k.Consumers

	var c *kafka.Consumer

	for _, consumer := range consumers {
		if kafkaConsumerName == consumer.Name {
			c = consumer.Consumer
		}
	}

	msg, err := c.ReadMessage(100)
	if err != nil {
		k.Log.Errorf("Encountered error when reading message. Error: %s", err)
		return "", "", err
	}

	return msg.String(), string(msg.Value), nil
}
