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
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := claims.(map[string]interface{})["email"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user by email
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req LetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedBody, err := utils.Encrypt(req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	letter := models.Letter{
		Email:     user.Email,
		BodyEnc:   encryptedBody,
		DeliverAt: req.DeliverAt,
		CreatedAt: time.Now(),
		Status:    "pending",
		UserID:    user.ID,
	}

	if err := initializers.DB.Create(&letter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save letter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Letter created successfully"})
}

func GetTodayLetters(c *gin.Context) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load timezone"})
		return
	}
	today := time.Now().In(loc).Format("2006-01-02")
	var ids []uint
	if err := initializers.DB.Model(&models.Letter{}).
		Where("DATE_FORMAT(deliver_at, '%Y-%m-%d') = ?", today).
		Pluck("id", &ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve letter ids"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ids": ids})
}

func GetLetterByEmail(c *gin.Context) {
	// Extract email from JWT claims (assuming middleware sets it in context)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no claims found"})
		return
	}

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
	ID      uint   `json:"id" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func SendEmails(c *gin.Context) {
	var req EmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var letter models.Letter
	if err := initializers.DB.First(&letter, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter not found"})
		return
	}

	decryptedBody, err := utils.Decrypt(letter.BodyEnc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt letter body"})
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Build HTML email body
	htmlBody := utils.BuildLetterHTML(letter, decryptedBody)
	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		fromName, from, letter.Email, req.Subject, htmlBody)

	cSMTP, err := smtp.Dial(addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to SMTP server"})
		return
	}
	defer cSMTP.Close()

	if err = cSMTP.Mail(from); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set sender"})
		return
	}
	if err = cSMTP.Rcpt(letter.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set recipient"})
		return
	}
	wc, err := cSMTP.Data()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data"})
		return
	}
	_, err = wc.Write([]byte(msg))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write message"})
		return
	}
	err = wc.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close writer"})
		return
	}

	// Update status to 'sent'
	letter.Status = "sent"
	initializers.DB.Save(&letter)
	c.JSON(http.StatusOK, gin.H{"message": "Email sent"})
}
