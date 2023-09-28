package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	HTTPAddress     string
	BaseURL         string
	FileStoragePath string
}

func InitConfig() *Config {
	config := &Config{}

	httpAddress := `:8080`
	baseURL := "http://localhost:8080/"
	fileStoragePath := "/tmp/short-url-db.json"
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		httpAddress = address
	}
	if url := os.Getenv("BASE_URL"); url != "" {
		baseURL = url
	}
	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		fileStoragePath = path
	}

	flag.StringVar(&config.HTTPAddress, "a", httpAddress, "HTTP Address")
	flag.StringVar(&config.BaseURL, "b", baseURL, "Base URL")
	flag.StringVar(&config.FileStoragePath, "f", fileStoragePath, "File storage path")
	flag.Parse()

	fmt.Printf("HTTP Address: %s\n", config.HTTPAddress)
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Printf("File storage path: %s\n", config.FileStoragePath)

	return config
}
