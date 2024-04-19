package coapcli

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Host          string `env:"COAP_CLI_HOST"`
	Port          string `env:"COAP_CLI_PORT"`
	ContentFormat int    `env:"COAP_CLI_CONTENT_FORMAT"`
	Auth          string `env:"COAP_CLI_AUTH"`
	Observe       bool   `env:"COAP_CLI_OBSERVE"`
}

func LoadConfig() (Config, error) {
	c := Config{}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	if err := env.Parse(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
