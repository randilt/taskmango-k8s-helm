package main

import (
	"log"
	"taskmango/apisvc/internal/config"
	"taskmango/apisvc/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup routes with auth middleware
	router := routes.SetupRouter(db, cfg)

	// Start main server
	go func() {
		log.Printf("API service started on port %d", cfg.APIPort)
		if err := router.Run(cfg.APIAddress()); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Start probe server
	probeRouter := routes.SetupProbeRouter(db)
	log.Printf("Health check service started on port %d", cfg.ProbePort)
	if err := probeRouter.Run(cfg.ProbeAddress()); err != nil {
		log.Fatalf("Failed to start probe server: %v", err)
	}
}