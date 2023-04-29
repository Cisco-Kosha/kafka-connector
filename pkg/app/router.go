package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	kpkg "github.com/kosha/kafka-connector/pkg/kafka"
	"net/http"
)

// getKafkaTopicMetadata godoc
// @Summary Get Kafka Topic Metadata details
// @Description Get Kafka Topic Metadata details
// @Tags topic
// @Accept json
// @Produce json
// @Success 200 {object} kafka.Metadata
// @Router /api/v1/getMetadata/{topic} [get]
func (a *App) getKafkaTopicMetadata(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	topic, _ := vars["topic"]

	kafkaTopicMetadata, err := a.Kafka.GetKafkaTopicMetadata(topic)
	if err != nil {
		a.Log.Errorf("Encountered an error in getKafkaTopicMetadata method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, &kafkaTopicMetadata)
}

// createKafkaTopic godoc
// @Summary Create Kafka Topic
// @Description Create new Kafka Topic
// @Tags topic
// @Accept json
// @Produce json
// @Param topic body kpkg.Topic false "Enter topic info"
// @Success 200 {object} []kafka.TopicResult
// @Router /api/v1/create [post]
func (a *App) createKafkaTopic(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var t kpkg.Topic
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload for topic creation")
		return
	}
	defer r.Body.Close()

	a.Log.Infof("Creating new kafka topic with name: %s if it does not exist already", t.Topic)

	results, err := a.Kafka.CreateTopic(t.Topic, t.NumParts, t.ReplicationFactor)
	if err != nil {
		a.Log.Errorf("Encountered an error in createTopics method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, results)
}

// deleteKafkaTopic godoc
// @Summary Delete Kafka Topic
// @Description Delete Kafka Topic
// @Tags topic
// @Accept json
// @Produce json
// @Param topic path string false "Enter topic name"
// @Success 200 {object} []kafka.TopicResult
// @Router /api/v1/delete/{topic} [delete]
func (a *App) deleteKafkaTopic(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	topic, _ := vars["topic"]

	results, err := a.Kafka.DeleteKafkaTopic(topic)
	if err != nil {
		a.Log.Errorf("Encountered an error in deleteKafkaTopic method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, &results)
}

// describeConfiguration godoc
// @Summary Describe Broker Config
// @Description Describe Broker Config
// @Tags metadata
// @Accept json
// @Produce json
// @Param resourceType path string false "Enter resource type"
// @Param resourceName path string false "Enter resource name"
// @Success 200 {object} []kafka.ConfigResourceResult
// @Router /api/v1/describe/{resourceType}/{resourceName} [get]
func (a *App) describeConfiguration(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	resourceType, _ := vars["resourceType"]
	resourceName, _ := vars["resourceName"]

	results, err := a.Kafka.DescribeConfig(resourceType, resourceName)
	if err != nil {
		a.Log.Errorf("Encountered an error in describeConfiguration method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, &results)
}

// produceKafkaMessage godoc
// @Summary Produce a kafka message
// @Description Produce a kafka message
// @Tags producer
// @Accept json
// @Produce json
// @Param topic path string false "Enter topic name"
// @Param message body map[string]interface{} false "Enter message"
// @Success 200 {object}  object
// @Router /api/v1/produce/{topic}/ [post]
func (a *App) produceKafkaMessage(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	topic, _ := vars["topic"]

	var u map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload for producing json kafka message")
		return
	}
	defer r.Body.Close()

	err := a.Kafka.ProduceMessages(topic, u)
	if err != nil {
		a.Log.Errorf("Encountered an error in createKafkaProducer method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// createKafkaConsumer godoc
// @Summary Create a new Kafka Consumer
// @Description Create a new Kafka Consumer
// @Tags consumer
// @Accept json
// @Produce json
// @Param resourceType path string false "Enter resource type"
// @Param resourceName path string false "Enter resource name"
// @Success 200 {object}  object
// @Router /api/v1/consumers [post]
func (a *App) createKafkaConsumer(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var c Consumers
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload for creating kafka consumers")
		return
	}
	defer r.Body.Close()
	fmt.Println("Printing the json body")
	a.Log.Infof("%v", c)
	err := a.Kafka.CreateConsumer(c.Name, c.Offset)
	if err != nil {
		a.Log.Errorf("Encountered an error in createKafkaConsumer method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success", "name": c.Name})
}

// subscribeToTopic godoc
// @Summary Subscribe to a Kafka Topic
// @Description Subscribe to a Kafka topic
// @Tags topic
// @Accept json
// @Produce json
// @Param name path string false "Enter consumer name"
// @Param topic body Topic false "Enter topic name"
// @Success 200 {object} object
// @Router /api/v1/consumers/{name}/subscribe [post]
func (a *App) subscribeToTopic(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	consumerName, _ := vars["name"]

	var t Topic
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		a.Log.Errorf("Error parsing json payload", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload for subscribing to topic")
		return
	}
	defer r.Body.Close()

	err := a.Kafka.SubscribeToTopics(t.Name, consumerName)
	if err != nil {
		a.Log.Errorf("Encountered an error in subscribeToTopic method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// consumeMessage godoc
// @Summary Consume last message from a kafka topic
// @Description Consume last message from a kafka topic
// @Tags consumer
// @Accept json
// @Produce json
// @Param name path string false "Enter consumer name"
// @Success 200 {object} map[string]string
// @Router /api/v1/consumers/{name}/consume [get]
func (a *App) consumeMessage(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	consumerName, _ := vars["name"]

	msg, err := a.Kafka.ConsumeMessage(consumerName)
	if err != nil {
		a.Log.Errorf("Encountered an error in consumeMessage method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var value map[string]interface{}
	err = json.Unmarshal(msg.Value, &value)
	if err != nil {
		a.Log.Errorf("Encountered an error in unmarshalling byte array into interface struct: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := &kpkg.ResMessage{
		TopicPartition: msg.TopicPartition,
		Value:          value,
		Key:            msg.Key,
		Timestamp:      msg.Timestamp,
		TimestampType:  msg.TimestampType,
		Opaque:         msg.Opaque,
		Headers:        msg.Headers,
	}

	respondWithJSON(w, http.StatusOK, res)
}
