package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/entity"
	"main.go/middleware"
	"main.go/service"
	"net/http"
	"strconv"
)

type TransactionsController struct {
	service service.TransactionsService
}

func NewTransactionsController(service service.TransactionsService) *TransactionsController {
	return &TransactionsController{service: service}
}

func (tc *TransactionsController) CreateTransaction(c *gin.Context) {
	claimUserID := c.GetUint("user_id")

	var transactionRequest entity.TransactionRequest
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if transactionRequest.UserID != claimUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: User ID mismatch"})
		return
	}

	transaction, err := tc.service.CreateTransaction(&transactionRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := service.ConvertToTransactionResponse(transaction)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully", "data": response})
}

func (tc *TransactionsController) GetTransactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := tc.service.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	claimUserID := c.GetUint("user_id")
	userRole := c.GetString("role")

	// 🔐 Validasi akses
	if userRole != "administrator" && transaction.UserID != claimUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	response := service.ConvertToTransactionResponse(transaction)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction fetched successfully", "data": response})
}

func (tc *TransactionsController) UpdateTransactionStatus(c *gin.Context) {
	middleware.Logger.Info("Controller: UpdateTransactionStatus called")

	if !middleware.HasRole(c, []string{"administrator"}) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid transaction ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var statusRequest entity.TransactionStatusRequest
	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		middleware.Logger.Error("Error binding request data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := tc.service.UpdateTransactionStatus(uint(id), statusRequest.Status); err != nil {
		middleware.Logger.Error("Failed to update transaction status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Transaction status updated successfully", zap.Uint("transaction_id", uint(id)))
	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}

func (tc *TransactionsController) DeleteTransaction(c *gin.Context) {
	middleware.Logger.Info("Controller: DeleteTransaction called")

	if !middleware.HasRole(c, []string{"administrator"}) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		middleware.Logger.Error("Invalid transaction ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	if err := tc.service.DeleteTransaction(uint(id)); err != nil {
		middleware.Logger.Error("Failed to delete transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Transaction deleted successfully", zap.Uint("transaction_id", uint(id)))
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func (tc *TransactionsController) GetAllTransactions(c *gin.Context) {
	middleware.Logger.Info("Controller: GetAllTransactions called")

	claimUserID := c.GetUint("user_id")
	userRole := c.GetString("role")

	var transactions []entity.Transaction
	var err error

	if userRole == "administrator" {
		// Jika administrator, panggil service untuk mendapatkan semua transaksi
		transactions, err = tc.service.GetAllTransactions()
		if err != nil {
			middleware.Logger.Error("Failed to fetch all transactions", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}
	} else {
		// Jika user biasa, panggil service untuk mendapatkan transaksi berdasarkan user_id
		transactions, err = tc.service.GetAllTransactionsByUser(claimUserID)
		if err != nil {
			middleware.Logger.Error("Failed to fetch user's transactions", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}
	}

	var responses []entity.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, service.ConvertToTransactionResponse(&transaction))
	}

	middleware.Logger.Info("Transactions fetched successfully", zap.Int("count", len(transactions)))
	c.JSON(http.StatusOK, gin.H{"message": "Transactions fetched successfully", "data": responses})
}

func (tc *TransactionsController) GetTransactionByUserID(c *gin.Context) {
	middleware.Logger.Info("Controller: GetTransactionByUserID called")

	// Ambil user_id dari parameter URL
	paramUserID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || paramUserID <= 0 {
		middleware.Logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Ambil user_id dari JWT token
	claimUserID := c.GetUint("user_id")
	userRole := c.GetString("role")

	// 🔐 Validasi akses
	if userRole != "administrator" && uint(paramUserID) != claimUserID {
		middleware.Logger.Warn("Access denied: User is not authorized to view other users' transactions")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Ambil transaksi berdasarkan user_id
	transactions, err := tc.service.GetAllTransactionsByUser(uint(paramUserID))
	if err != nil {
		middleware.Logger.Error("Failed to fetch transactions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// Kirimkan respons ke client
	middleware.Logger.Info("Transactions fetched successfully", zap.Int("count", len(transactions)))
	c.JSON(http.StatusOK, gin.H{"message": "Transactions fetched successfully", "data": transactions})
}
