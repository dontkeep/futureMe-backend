package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	Letters      []Letter `gorm:"foreignKey:UserID"`
}

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
