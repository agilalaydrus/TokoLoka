package entity

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FullName    string    `gorm:"not null" json:"full_name"`
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Email       string    `gorm:"unique" json:"email,omitempty"` // Email opsional
	Password    string    `gorm:"type:varchar(255);not null" json:"password"`
	Address     string    `json:"address"`
	Role        string    `gorm:"size:20;not null" json:"role"` // user / administrator
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
