package models

import (
	"time"
)

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"unique;not null"`
	PasswordHash     string `gorm:"not null"`
	CreatedAt        time.Time
	ResetToken       string `gorm:"size:64"`
	ResetTokenExpiry time.Time
	Letters          []Letter `gorm:"foreignKey:UserID"`
}
