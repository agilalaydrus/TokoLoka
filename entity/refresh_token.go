package entity

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}
