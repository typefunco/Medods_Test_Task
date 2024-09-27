package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	DatabaseURL   string `json:"database_url"`
	JWTSecret     string `json:"jwt_secret"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Error decoding config file: %v", err)
		return nil, err
	}

	fmt.Printf("Config successfully loaded: %+v", config)
	return &config, nil
}
