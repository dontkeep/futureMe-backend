package utils

import (
	"dontkeep/futureme-backend/initializers"
	"dontkeep/futureme-backend/models"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func StartLetterScheduler() {
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Minute().Do(SendScheduledEmails)
	s.StartAsync()
}

func SendScheduledEmails() {
	var letters []models.Letter
	if err := initializers.DB.Where("deliver_at <= ? AND status = ?", time.Now(), "pending").Find(&letters).Error; err != nil {
		fmt.Println("Error fetching letters:", err)
		return
	}
	for _, letter := range letters {
		decryptedBody, err := Decrypt(letter.BodyEnc)
		if err != nil {
			fmt.Println("Decryption failed:", err)
			continue
		}
		htmlBody := BuildLetterHTML(letter, decryptedBody)
		err = SendEmail(letter.Email, "Your Future Letter", htmlBody)
		if err != nil {
			fmt.Println("Failed to send email:", err)
			continue
		}
		letter.Status = "sent"
		now := time.Now()
		letter.SentAt = &now
		if err := initializers.DB.Save(&letter).Error; err != nil {
			fmt.Println("Failed to update status:", err)
		} else {
			fmt.Printf("Email sent successfully to %s\n", letter.Email)
		}
	}
}
