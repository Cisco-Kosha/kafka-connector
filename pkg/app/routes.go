package app

import (
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) initializeRoutes() {
	var apiV1 = "/api/v1"
	a.Router.HandleFunc(apiV1+"/create", a.createKafkaTopic).Methods("POST", "OPTIONS").Name("CreateTopicHandler")
	a.Router.HandleFunc(apiV1+"/produce/{topic}", a.produceKafkaMessage).Methods("POST", "OPTIONS").Name("ProduceMessageToTopicHandler")
	a.Router.HandleFunc(apiV1+"/delete/{topic}", a.deleteKafkaTopic).Methods("DELETE", "OPTIONS").Name("DeleteTopicHandler")
	a.Router.HandleFunc(apiV1+"/describe/{resourceType}/{resourceName}", a.describeConfiguration).Methods("GET", "OPTIONS").Name("DescribeConfiguration")
	a.Router.HandleFunc(apiV1+"/getMetadata/{topic}", a.getKafkaTopicMetadata).Methods("GET", "OPTIONS").Name("GetKafkaTopicMetadataHandler")
	a.Router.HandleFunc(apiV1+"/consumers", a.createKafkaConsumer).Methods("POST", "OPTIONS").Name("CreateKafkaConsumer")
	a.Router.HandleFunc(apiV1+"/consumers/{name}/subscribe", a.subscribeToTopic).Methods("POST", "OPTIONS").Name("SubscribeToTopicForKafkaConsumer")
	a.Router.HandleFunc(apiV1+"/consumers/{name}/consume", a.consumeMessage).Methods("GET", "OPTIONS").Name("ConsumeMessageForKafkaConsumer")

	a.Router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

}
