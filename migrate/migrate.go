package main

import (
	"dontkeep/futureme-backend/initializers"
	"dontkeep/futureme-backend/models"
	"log"
)

func main() {
	initializers.ConnectDB()
	if err := initializers.DB.AutoMigrate(&models.User{}, &models.Letter{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Migration completed successfully.")
}
