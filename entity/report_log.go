package entity

import "time"

type ReportLog struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	ReportName  string `gorm:"size:100;not null"`
	Filters     string `gorm:"type:text"`
	GeneratedAt time.Time
	User        User `gorm:"foreignKey:UserID"`
}
