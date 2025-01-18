package repository

import (
	"gorm.io/gorm"
	"main.go/entity"
)

type TransactionsRepository interface {
	Create(transaction *entity.Transaction) error
	GetByID(id uint) (*entity.Transaction, error)
	GetAllByUserID(userID uint) ([]entity.Transaction, error)
	GetAll() ([]entity.Transaction, error) // ✅ Tambahkan method ini
	Update(transaction *entity.Transaction) error
	Delete(id uint) error
}

type transactionsRepository struct {
	db *gorm.DB
}

func NewTransactionsRepository(db *gorm.DB) TransactionsRepository {
	return &transactionsRepository{db: db}
}

// ✅ Create - Membuat transaksi baru
func (r *transactionsRepository) Create(transaction *entity.Transaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

// ✅ GetByID - Mengambil transaksi berdasarkan ID
func (r *transactionsRepository) GetByID(id uint) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, email, role")
	}).Preload("Items.Product.Category").First(&transaction, id).Error

	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ✅ GetAllByUserID - Mengambil semua transaksi berdasarkan User ID
func (r *transactionsRepository) GetAllByUserID(userID uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.
		Preload("User").
		Preload("Items.Product.Category").
		Where("user_id = ?", userID).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionsRepository) GetAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Preload("User").Preload("Items.Product.Category").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// ✅ Update - Mengupdate transaksi
func (r *transactionsRepository) Update(transaction *entity.Transaction) error {
	if err := r.db.Save(transaction).Error; err != nil {
		return err
	}
	return nil
}

// ✅ Delete - Menghapus transaksi berdasarkan ID
func (r *transactionsRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.Transaction{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionsRepository) CreateActivityLog(log *entity.ActivityLog) error {
	return r.db.Create(log).Error
}
