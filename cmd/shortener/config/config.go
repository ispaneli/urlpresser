package config

import "flag"

type Config struct {
	HTTPAddress string
	BaseURL     string
}

func InitConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.HTTPAddress, "a", `:8080`, "HTTP Address")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8000/", "Base URL")
	flag.Parse()

	return config
}
