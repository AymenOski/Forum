package main

import (
	"fmt"
	"log"

	"forum/infrastructure/database"
	"forum/infrastructure/server"
)

func main() {
		fmt.Println("Server started on http://localhost:8080")
	if err := server.Froum_server().ListenAndServe(); err != nil {
		log.Fatalf("500 - Internal Server Error: %v", err)
	}

	db := database.SetingUpDB()
	defer db.Close()
}
