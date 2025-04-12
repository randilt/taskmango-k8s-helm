package main

import (
	"log"
	"taskmango/authsvc/internal/config"
	"taskmango/authsvc/internal/routes"
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

	// Setup routes
	router := routes.SetupRouter(db, cfg)

	// Start main server
	go func() {
		log.Printf("Auth service started on port %d", cfg.AuthPort)
		if err := router.Run(cfg.AuthAddress()); err != nil {
			log.Fatalf("Failed to start auth server: %v", err)
		}
	}()

	// Start probe server
	probeRouter := routes.SetupProbeRouter(db)
	log.Printf("Health check service started on port %d", cfg.ProbePort)
	if err := probeRouter.Run(cfg.ProbeAddress()); err != nil {
		log.Fatalf("Failed to start probe server: %v", err)
	}
}