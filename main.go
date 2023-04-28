package main

import (
	"fmt"
	"github.com/kosha/kafka-connector/pkg/app"
	logger "github.com/kosha/kafka-connector/pkg/logger"
)

const (
	port = 8000
)

func main() {
	a := app.App{}

	log := logger.New("app", "kafka-connector")

	a.Initialize(log)

	log.Infof("Starting Kafka Connector on port: %d", port)

	a.Run(fmt.Sprintf(":%d", port))
}
