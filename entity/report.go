package entity

import "time"

type ReportLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                 // Primary Key
	UserID      uint      `gorm:"not null;index" json:"user_id"`        // User ID dengan index
	ReportName  string    `gorm:"size:100;not null" json:"report_name"` // Nama laporan
	Filters     string    `gorm:"type:text" json:"filters"`             // Filter dalam bentuk string (bisa JSON)
	GeneratedAt time.Time `gorm:"autoCreateTime" json:"generated_at"`   // Waktu pembuatan laporan
	User        User      `gorm:"foreignKey:UserID" json:"user"`        // Relasi ke entitas User
}

type ReportService interface {
	GetTransactionSummaryAll(filters ReportFilters) ([]TransactionSummary, error)
	GetTransactionSummaryByUser(userID uint, filters ReportFilters) ([]TransactionSummary, error)
	GetTransactionSummaryByProduct(productID uint, filters ReportFilters) ([]TransactionSummary, error)
	GetTransactionSummaryByCategory(categoryID uint, filters ReportFilters) ([]TransactionSummary, error)
	ExportReportToCSV(data []TransactionSummary) ([]byte, error)
	ExportReportToExcel(data []TransactionSummary) ([]byte, error)
}

type TransactionSummary struct {
	TransactionID   uint    `json:"transaction_id"`
	UserID          uint    `json:"user_id"`   // Digunakan di repository
	UserName        string  `json:"user_name"` // Digunakan di repository
	ProductName     string  `json:"product_name"`
	CategoryName    string  `json:"category_name"`
	Quantity        int     `json:"quantity"`
	TotalPrice      float64 `json:"total_price"`
	TransactionDate string  `json:"transaction_date"`
}

type ReportFilters struct {
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	Username    string `json:"username,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	UserID      uint   `json:"-"`
	Page        int    `json:"page,omitempty"`  // Halaman saat ini
	Limit       int    `json:"limit,omitempty"` // Jumlah data per halaman
}
