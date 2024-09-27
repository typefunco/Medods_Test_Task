package main

import (
	"log"
	"medods_auth/authService/internal/app"
	"medods_auth/authService/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("/root/config/cfg.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app.RunServer(cfg)
}
