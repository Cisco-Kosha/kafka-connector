package app

func (a *App) initializeRoutes() {
	var apiV1 = "/api/v1"
	a.Router.HandleFunc(apiV1+"/create", a.createKafkaTopic).Methods("POST").Name("CreateTopicHandler")
	a.Router.HandleFunc(apiV1+"/produce/{topic}", a.produceKafkaMessage).Methods("POST").Name("ProduceMessageToTopicHandler")
	a.Router.HandleFunc(apiV1+"/delete/{topic}", a.deleteKafkaTopic).Methods("DELETE").Name("DeleteTopicHandler")
	a.Router.HandleFunc(apiV1+"/describe/{resourceType}/{resourceName}", a.describeConfiguration).Methods("GET").Name("DescribeConfiguration")
	a.Router.HandleFunc(apiV1+"/getMetadata/{topic}", a.getKafkaTopicMetadata).Methods("GET").Name("GetKafkaTopicMetadataHandler")
	a.Router.HandleFunc(apiV1+"/consumers", a.createKafkaConsumer).Methods("POST").Name("CreateKafkaConsumer")
	a.Router.HandleFunc(apiV1+"/consumers/{name}/subscribe", a.subscribeToTopic).Methods("POST").Name("SubscribeToTopicForKafkaConsumer")
	a.Router.HandleFunc(apiV1+"/consumers/{name}/consume", a.consumeMessage).Methods("GET").Name("ConsumeMessageForKafkaConsumer")

}
