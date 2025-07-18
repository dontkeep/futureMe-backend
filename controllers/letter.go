package controllers

import (
	"dontkeep/futureme-backend/initializers"
	"dontkeep/futureme-backend/models"
	"dontkeep/futureme-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LetterRequest struct {
	Email     string    `json:"email" binding:"required,email"`
	Body      string    `json:"body" binding:"required"`
	DeliverAt time.Time `json:"deliverAt" binding:"required"`
}

func CreateLetter(c *gin.Context) {
	var req LetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedBody, err := utils.Encrypt(req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt letter body"})
		return
	}

	letter := models.Letter{
		Email:     req.Email,
		BodyEnc:   encryptedBody,
		DeliverAt: req.DeliverAt,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	if err := initializers.DB.Create(&letter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save letter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Letter saved successfully!"})
}
