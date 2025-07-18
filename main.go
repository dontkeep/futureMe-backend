package main

import (
	"dontkeep/futureme-backend/controllers"
	"dontkeep/futureme-backend/initializers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Duplicate import block removed

func main() {
	// Initialize logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	log.Logger = logger

	// Connect to DB
	initializers.ConnectDB()

	// Setup Gin router
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is working"})
	})
	r.POST("/letters", controllers.CreateLetter)

	log.Info().Msgf("Starting server on port %d", 3000)
	if err := r.Run(":3000"); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
