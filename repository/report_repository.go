package repository

import (
	"gorm.io/gorm"
	"main.go/entity"
)

type ReportRepository interface {
	GetTransactionsSummary(filters entity.ReportFilters, isAdmin bool, userID uint) ([]entity.TransactionSummary, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetTransactionsSummary(filters entity.ReportFilters, isAdmin bool, userID uint) ([]entity.TransactionSummary, error) {
	var summaries []entity.TransactionSummary

	query := r.db.Table("transactions").
		Select(`
			transactions.id as transaction_id,
			users.id as user_id,
			users.username as user_name,
			products.name as product_name,
			categories.name as category_name,
			sum(transaction_items.quantity) as quantity,
			sum(transaction_items.quantity * transaction_items.price) as total_price,
			transactions.created_at as transaction_date
		`).
		Joins("join transaction_items on transactions.id = transaction_items.transaction_id").
		Joins("join products on transaction_items.product_id = products.id").
		Joins("join categories on products.category_id = categories.id").
		Joins("join users on transactions.user_id = users.id").
		Group("transactions.id, users.id, products.name, categories.name, transactions.created_at")

	if filters.StartDate != "" && filters.EndDate != "" {
		query = query.Where("transactions.created_at BETWEEN ? AND ?", filters.StartDate, filters.EndDate)
	}

	if !isAdmin {
		query = query.Where("transactions.user_id = ?", userID)
	}

	if filters.Username != "" {
		query = query.Where("users.username LIKE ?", "%"+filters.Username+"%")
	}

	if filters.ProductName != "" {
		query = query.Where("products.name LIKE ?", "%"+filters.ProductName+"%")
	}

	// Pagination
	offset := (filters.Page - 1) * filters.Limit
	query = query.Limit(filters.Limit).Offset(offset)

	if err := query.Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}
