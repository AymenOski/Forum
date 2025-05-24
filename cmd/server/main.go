package main

import (
	"fmt"
	"log"

	"forum/config"
		"forum/infrastructure/database"
	"forum/infrastructure/server"
)

func main() {
	cfg := config.Load()
	db := database.SetingUpDB(cfg.DatabasePath)
	defer db.Close()

	fmt.Println("Server started on http://localhost:8080")
	if err := server.Froum_server().ListenAndServe(); err != nil {
		log.Fatalf("500 - Internal Server Error: %v", err)
	}
}
