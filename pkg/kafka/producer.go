package kafka

import (
	"bytes"
	"encoding/gob"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateProducer(brokers string) *kafka.Producer {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})

	if err != nil {
		return nil
	}

	return p
}

func (k *Kafka) ProduceMessages(topic string, value map[string]interface{}) error {

	deliveryChan := make(chan kafka.Event)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		k.Log.Errorf("Cannot convert payload into byte array using gob package: %v", value)
		return err
	}
	err = k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          buf.Bytes(),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		k.Log.Errorf("Delivery failed: %v\n", m.TopicPartition.Error)
		return err
	} else {
		k.Log.Infof("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil

}
