package entity

import "time"

// Transaction struct untuk merepresentasikan transaksi
type Transaction struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	UserID            uint              `gorm:"not null" json:"user_id"`
	DestinationNumber string            `gorm:"size:15" json:"destination_number"`
	TotalPrice        float64           `gorm:"type:decimal(10,2)" json:"total_price"`
	Status            string            `gorm:"size:20;default:'pending'" json:"status"`
	SerialNumber      string            `gorm:"size:50" json:"serial_number"` // Tambahkan ini
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	User              User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items             []TransactionItem `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
}

// TransactionItem struct untuk merepresentasikan item dalam transaksi
type TransactionItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null" json:"transaction_id"`
	ProductID     uint      `gorm:"not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Product       Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TransactionRequest struct untuk menerima request transaksi dari client
type TransactionRequest struct {
	UserID            uint                     `json:"user_id"`
	DestinationNumber string                   `json:"destination_number"` // Nomor tujuan transaksi
	Items             []TransactionItemRequest `json:"items"`
}

// TransactionItemRequest struct untuk menerima item dalam request transaksi
type TransactionItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// TransactionStatusRequest struct untuk menerima request perubahan status transaksi
type TransactionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// TransactionCallbackResponse struct untuk mengirimkan callback response dari supplier
type TransactionCallbackResponse struct {
	RequestID         uint      `json:"request_id"`         // ID transaksi dari TokoLoka
	SerialNumber      string    `json:"serial_number"`      // Nomor seri unik dari supplier
	Status            string    `json:"status"`             // Status transaksi (success/failed)
	DestinationNumber string    `json:"destination_number"` // Nomor tujuan transaksi
	TotalPrice        float64   `json:"total_price"`        // Total harga transaksi
	Message           string    `json:"message"`            // Pesan dari supplier
	CallbackTime      time.Time `json:"callback_time"`      // Waktu callback diterima
}
