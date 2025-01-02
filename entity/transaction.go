package entity

import "time"

type Transaction struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null"`
	TotalPrice float64   `gorm:"not null"`
	Status     string    `gorm:"default:pending;not null"`
	CreatedAt  time.Time `gorm:"type:timestamp(3);autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp(3);autoUpdateTime"`
	User       User      `gorm:"foreignKey:UserID"`
}

type TransactionItem struct {
	ID            uint        `gorm:"primaryKey;autoIncrement"`
	TransactionID uint        `gorm:"not null"`
	ProductID     uint        `gorm:"not null"`
	Quantity      int         `gorm:"not null"`
	Price         float64     `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"type:timestamp(3);autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"type:timestamp(3);autoUpdateTime"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID"`
	Product       Product     `gorm:"foreignKey:ProductID"`
}
