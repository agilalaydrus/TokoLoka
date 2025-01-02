package entity

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	CreatedAt   time.Time `gorm:"type:timestamp(3);autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp(3);autoUpdateTime" json:"updated_at"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp(3);autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp(3);autoUpdateTime" json:"updated_at"`
}
