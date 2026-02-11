package main

import (
	"english-learning/configs"
	"english-learning/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env (ignore error if not present)
	_ = godotenv.Load()

	// Load Config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init Application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// Run Server
	if err := application.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
