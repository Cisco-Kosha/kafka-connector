package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kosha/kafka-connector/pkg/config"
	logger "github.com/kosha/kafka-connector/pkg/logger"

	kpkg "github.com/kosha/kafka-connector/pkg/kafka"
)

// App represents the application
type App struct {
	Router *mux.Router
	Cfg    *config.Config
	Log    logger.Logger
	Kafka  *kpkg.Kafka
}

type Topic struct {
	Name string `json:"name,omitempty"`
}

type Consumers struct {
	Name   string `json:"name,omitempty"`
	Offset string `json:"offset,omitempty"`
}

func router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

// Initialize sets up the db connection and routes for the app
func (a *App) Initialize(log logger.Logger) {

	cfg := config.Get()

	a.Router = router()

	a.Log = log
	a.Kafka = kpkg.NewKafkaClient(cfg, log)

	a.initializeRoutes()
}

// Run starts the app and serves on the specified addr
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
