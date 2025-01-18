package entity

import "time"

// TransactionResponse - Struct untuk response transaksi
type TransactionResponse struct {
	ID                uint                      `json:"id"`
	UserID            uint                      `json:"user_id"`
	DestinationNumber string                    `json:"destination_number"`
	TotalPrice        float64                   `json:"total_price"`
	Status            string                    `json:"status"`
	SerialNumber      string                    `json:"serial_number"` // Tambahkan ini
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	User              UserSafeResponse          `json:"user"`
	Items             []TransactionItemResponse `json:"items"`
}

type UserSafeResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type TransactionItemResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"`
	Product   ProductResponse `json:"product"`
}

type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	Category    CategoryResponse `json:"category"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
