package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"main.go/entity"
	"main.go/middleware"
	"main.go/repository"
	"math/rand"
	"time"
)

type TransactionsService interface {
	CreateTransaction(transactionRequest *entity.TransactionRequest) (*entity.Transaction, error)
	GetAllTransactions() ([]entity.Transaction, error)
	GetTransactionByID(id uint) (*entity.Transaction, error)
	GetAllTransactionsByUser(userID uint) ([]entity.Transaction, error)
	UpdateTransactionStatus(id uint, status string) error
	DeleteTransaction(id uint) error
}

type transactionsService struct {
	repository         repository.TransactionsRepository
	productRepo        repository.ProductRepository
	activityLogService ActivityLogService
}

func NewTransactionsService(repo repository.TransactionsRepository, productRepo repository.ProductRepository, activityLogService ActivityLogService) TransactionsService {
	return &transactionsService{
		repository:         repo,
		productRepo:        productRepo,
		activityLogService: activityLogService,
	}
}

// CreateTransaction - Membuat transaksi baru
func (s *transactionsService) CreateTransaction(transactionRequest *entity.TransactionRequest) (*entity.Transaction, error) {
	middleware.Logger.Info("Service: CreateTransaction called")

	// Validasi Destination Number
	if len(transactionRequest.DestinationNumber) < 11 || len(transactionRequest.DestinationNumber) > 12 {
		middleware.Logger.Warn("Invalid destination number")
		return nil, errors.New("invalid destination number")
	}

	// Proses transaksi
	transaction := &entity.Transaction{
		UserID:            transactionRequest.UserID,
		DestinationNumber: transactionRequest.DestinationNumber,
		Status:            "pending",
	}

	// Hitung total harga berdasarkan produk di database
	totalPrice := 0.0
	for _, item := range transactionRequest.Items {
		// Ambil harga produk dari database
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			middleware.Logger.Error("Product not found", zap.Uint("product_id", item.ProductID))
			return nil, errors.New("product not found")
		}

		// Hitung total harga untuk item ini
		itemTotalPrice := product.Price * float64(item.Quantity)

		// Tambahkan ke total transaksi
		totalPrice += itemTotalPrice

		// Simpan item transaksi
		transaction.Items = append(transaction.Items, entity.TransactionItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	// Tetapkan total harga yang dihitung
	transaction.TotalPrice = totalPrice

	// Simpan transaksi ke database
	if err := s.repository.Create(transaction); err != nil {
		middleware.Logger.Error("Failed to create transaction", zap.Error(err))
		return nil, err
	}

	// Simulasi callback
	go s.simulateCallback(transaction)

	return s.GetTransactionByID(transaction.ID)
}

func (s *transactionsService) simulateCallback(transaction *entity.Transaction) {
	time.Sleep(2 * time.Second) // Simulasi waktu tunggu supplier

	randomSerial := s.generateSerialNumber()
	isFailed := false
	var failReason string

	// Validasi berdasarkan setiap item dalam transaksi
	calculatedTotal := 0.0
	for _, item := range transaction.Items {
		calculatedTotal += item.Price * float64(item.Quantity)
	}

	// Validasi jika TotalPrice tidak sesuai dengan perhitungan
	if calculatedTotal != transaction.TotalPrice {
		isFailed = true
		failReason = "Total price mismatch"
	}

	// Simulasi callback sukses/gagal
	if isFailed {
		middleware.Logger.Warn("Transaction failed",
			zap.Uint("transaction_id", transaction.ID),
			zap.String("reason", failReason),
		)
		transaction.Status = "failed"
	} else {
		middleware.Logger.Info("Transaction success",
			zap.String("destination_number", transaction.DestinationNumber),
			zap.Uint("transaction_id", transaction.ID),
			zap.String("serial_number", randomSerial),
		)
		transaction.Status = "success"
		transaction.SerialNumber = randomSerial
	}

	// Update status transaksi di database
	if err := s.repository.Update(transaction); err != nil {
		middleware.Logger.Error("Failed to update transaction status", zap.Error(err))
	}

	// Log aktivitas callback
	callbackDetails := fmt.Sprintf("Transaction ID: %d, Status: %s, Serial Number: %s", transaction.ID, transaction.Status, transaction.SerialNumber)
	s.activityLogService.CreateActivityLog(transaction.UserID, "Callback Received", callbackDetails)
}

func (s *transactionsService) generateSerialNumber() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	serial := make([]byte, 8)
	for i := range serial {
		serial[i] = charset[rand.Intn(len(charset))]
	}

	return fmt.Sprintf("SN-%s", string(serial))
}

// GetTransactionByID - Mengambil detail transaksi berdasarkan ID
func (s *transactionsService) GetTransactionByID(id uint) (*entity.Transaction, error) {
	middleware.Logger.Info("Service: GetTransactionByID called", zap.Uint("transaction_id", id))

	transaction, err := s.repository.GetByID(id)
	if err != nil {
		middleware.Logger.Error("Service: Transaction not found", zap.Uint("transaction_id", id), zap.Error(err))
		return nil, errors.New("transaction not found")
	}

	middleware.Logger.Info("Service: Transaction fetched successfully", zap.Any("transaction", transaction))
	return transaction, nil
}

// GetAllTransactionsByUser - Mendapatkan semua transaksi milik user
func (s *transactionsService) GetAllTransactionsByUser(userID uint) ([]entity.Transaction, error) {
	middleware.Logger.Info("Service: GetAllTransactionsByUser called", zap.Uint("user_id", userID))

	transactions, err := s.repository.GetAllByUserID(userID)
	if err != nil {
		middleware.Logger.Error("Service: Error fetching transactions", zap.Error(err))
		return nil, err
	}

	middleware.Logger.Info("Service: Transactions fetched successfully", zap.Int("count", len(transactions)))
	return transactions, nil
}

// GetAllTransactions - Mendapatkan semua transaksi (khusus untuk admin)
func (s *transactionsService) GetAllTransactions() ([]entity.Transaction, error) {
	middleware.Logger.Info("Service: GetAllTransactions called")

	transactions, err := s.repository.GetAll()
	if err != nil {
		middleware.Logger.Error("Service: Error fetching all transactions", zap.Error(err))
		return nil, err
	}

	middleware.Logger.Info("Service: All transactions fetched successfully", zap.Int("count", len(transactions)))
	return transactions, nil
}

// UpdateTransactionStatus - Mengupdate status transaksi
func (s *transactionsService) UpdateTransactionStatus(id uint, status string) error {
	middleware.Logger.Info("Service: UpdateTransactionStatus called", zap.Uint("transaction_id", id), zap.String("status", status))

	validStatuses := map[string]bool{
		"pending": true,
		"process": true,
		"failed":  true,
		"success": true,
	}

	if !validStatuses[status] {
		middleware.Logger.Warn("Service: Invalid transaction status", zap.String("status", status))
		return errors.New("invalid transaction status")
	}

	transaction, err := s.repository.GetByID(id)
	if err != nil {
		middleware.Logger.Error("Service: Transaction not found", zap.Uint("transaction_id", id), zap.Error(err))
		return errors.New("transaction not found")
	}

	transaction.Status = status

	if err := s.repository.Update(transaction); err != nil {
		middleware.Logger.Error("Service: Failed to update transaction status", zap.Error(err))
		return errors.New("failed to update transaction status")
	}

	middleware.Logger.Info("Service: Transaction status updated successfully", zap.Uint("transaction_id", transaction.ID))
	return nil
}

// DeleteTransaction - Menghapus transaksi berdasarkan ID
func (s *transactionsService) DeleteTransaction(id uint) error {
	middleware.Logger.Info("Service: DeleteTransaction called", zap.Uint("transaction_id", id))

	_, err := s.repository.GetByID(id)
	if err != nil {
		middleware.Logger.Error("Service: Transaction not found for deletion", zap.Uint("transaction_id", id))
		return errors.New("transaction not found")
	}

	if err := s.repository.Delete(id); err != nil {
		middleware.Logger.Error("Service: Failed to delete transaction", zap.Error(err))
		return errors.New("failed to delete transaction")
	}

	middleware.Logger.Info("Service: Transaction deleted successfully", zap.Uint("transaction_id", id))
	return nil
}
