package models

import (
	"time"
)

type Letter struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null"`
	BodyEnc   string    `gorm:"not null"`
	DeliverAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	SentAt    *time.Time
	Status    string `gorm:"default:'pending'"`
	UserID    *uint
}
