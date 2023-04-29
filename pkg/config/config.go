package config

import (
	"flag"
	"os"
)

type Config struct {
	brokerAddress string
	username      string
	password      string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.brokerAddress, "brokerAddress", os.Getenv("BROKER_ADDRESS"), "Broker Address")
	flag.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Broker Username")
	flag.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Broker Password")

	flag.Parse()

	return conf
}

func (c *Config) GetBrokerAddress() string {
	return c.brokerAddress
}

func (c *Config) GetBrokerUsername() string {
	return c.username
}

func (c *Config) GetBrokerPassword() string {
	return c.password
}
