package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"unique;not null"`
	Role      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp(3);autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp(3);autoUpdateTime"`
}
