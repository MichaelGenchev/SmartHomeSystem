package main

import (
	"context"
	"log"
	"os"

	"github.com/MichaelGenchev/smart-home-system/internal/app/userapp"
	"github.com/MichaelGenchev/smart-home-system/internal/common/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and initialize the device service application
	app := deviceapp.NewApp(cfg)
	if err := app.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize device service: %v", err)
	}

	// Run the device service
	if err := app.Run(); err != nil {
		log.Printf("Device service error: %v", err)
		os.Exit(1)
	}
}