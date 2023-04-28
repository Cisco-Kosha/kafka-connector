package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	kpkg "github.com/kosha/kafka-connector/pkg/kafka"
	"net/http"
)

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

func (a *App) consumeMessage(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	consumerName, _ := vars["name"]

	topicMetadata, msg, err := a.Kafka.ConsumeMessage(consumerName)
	if err != nil {
		a.Log.Errorf("Encountered an error in consumeMessage method: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"KafkaMessage": topicMetadata, "value": msg})
}
