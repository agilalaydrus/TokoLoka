package entity

import (
	"time"
)

type ActivityLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"` // Contoh: "Create Transaction", "Callback Received"
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
