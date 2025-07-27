package controllers

import (
	"dontkeep/futureme-backend/initializers"
	"dontkeep/futureme-backend/models"
	"dontkeep/futureme-backend/utils"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
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

func GetTodayLetters(c *gin.Context) {
	var letters []models.Letter
	today := time.Now().Format("2006-01-02")
	if err := initializers.DB.Where("DATE(deliver_at) = ?", today).Find(&letters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve letters"})
		return
	}

	for i, letter := range letters {
		decryptedBody, err := utils.Decrypt(letter.BodyEnc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt letter body"})
			return
		}
		letters[i].BodyEnc = decryptedBody
	}

	c.JSON(http.StatusOK, letters)
}

func GetLetterByEmail(c *gin.Context) {
	// Extract email from JWT claims (assuming middleware sets it in context)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no claims found"})
		return
	}

	// You may need to adjust this depending on your JWT middleware implementation
	email, ok := claims.(map[string]interface{})["email"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: email not found in token"})
		return
	}

	var letters []models.Letter
	if err := initializers.DB.Where("email = ?", email).Find(&letters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve letters"})
		return
	}

	for i, letter := range letters {
		decryptedBody, err := utils.Decrypt(letter.BodyEnc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt letter body"})
			return
		}
		letters[i].BodyEnc = decryptedBody
	}

	c.JSON(http.StatusOK, letters)
}

type EmailSendRequest struct {
	To      []string `json:"to" binding:"required"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
}

func SendEmails(c *gin.Context) {
	var req EmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	for _, to := range req.To {
		msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			fromName, from, to, req.Subject, req.Body)

		// No auth for local postfix relay
		err := smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send to %s: %v", to, err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent"})
}
