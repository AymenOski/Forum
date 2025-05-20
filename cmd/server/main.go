package main

import (
	"forum/infrastructure/database"
)

func main() {
	db := database.SetingUpDB()
	defer db.Close()
}
