package models

import "time"

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	Token     string
	UserID    uint
	ExpiresAt time.Time
}
