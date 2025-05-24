package main

import (
	"fmt"
	"log"

	"forum/infrastructure/database"
	"forum/infrastructure/server"
)

func main() {
	db := database.SetingUpDB()
	defer db.Close()
	fmt.Println("Server started on http://localhost:8080")
	if err := server.Froum_server().ListenAndServe(); err != nil {
		log.Fatalf("500 - Internal Server Error: %v", err)
	}
}
