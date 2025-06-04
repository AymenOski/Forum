package main

import (
	"fmt"
	"html/template"
	"log"

	"forum/config"
	"forum/infrastructure/database"
	"forum/infrastructure/server"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Printf("Warning: Failed to initialize templates: %v", err)
	}
}

func main() {
	cfg := config.Load()
	db := database.SetingUpDB(cfg.DatabasePath)
	defer db.Close()

	fmt.Println("Server started on http://localhost:8080")
	if err := server.MyServer(db, templates).ListenAndServe(); err != nil {
		log.Fatalf("500 - Internal Server Error: %v", err)
	}
}
