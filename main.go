package main

import (
	"dontkeep/futureme-backend/controllers"
	"dontkeep/futureme-backend/initializers"
	"dontkeep/futureme-backend/middleware"
	"dontkeep/futureme-backend/utils"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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
	utils.StartLetterScheduler()

	// Setup Gin router
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is working"})
	})
	r.POST("/letters", middleware.RequireAuth(), middleware.RateLimit(), controllers.CreateLetter)
	r.GET("/letters/today", middleware.RequireStaticToken(), controllers.GetTodayLetters)
	r.GET("/letters", middleware.RequireAuth(), controllers.GetLetterByEmail)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/forgot-password", controllers.ForgotPassword)
	r.POST("/reset-password", controllers.ResetPassword)

	log.Info().Msgf("Starting server on port %d", 3000)
	if err := r.Run(":3000"); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
