package main

import (
	"fmt"
	"log"
	"os"

	"forum/config"
	"forum/infrastructure/database"
	"forum/infrastructure/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database
	db := database.SetingUpDB(cfg.DatabasePath)
	defer db.Close()

	// Create server
	srv := server.Forum_server(db)

	// Get port from environment or config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Set custom port if provided
	srv.SetPort(port)

	// Start server
	fmt.Printf("Server started on http://localhost:%s\n", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("500 - Internal Server Error: %v", err)
	}
}
