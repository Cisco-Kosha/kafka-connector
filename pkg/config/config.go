package config

import (
	"flag"
	"os"
)

type Config struct {
	brokerAddress string
	username      string
	password      string
	privKey       string
	pubKey        string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.brokerAddress, "brokerAddress", os.Getenv("BROKER_ADDRESS"), "Broker Address")
	flag.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Broker Username")
	flag.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Broker Password")
	flag.StringVar(&conf.privKey, "privateKey", os.Getenv("PRIVATE_KEY_PATH"), "Private Key Path")
	flag.StringVar(&conf.pubKey, "publicKey", os.Getenv("PUBLIC_KEY_PATH"), "Public Key Path")

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

func (c *Config) GetPrivKeyPath() string {
	return c.privKey
}

func (c *Config) GetPubKeyPath() string {
	return c.pubKey
}
