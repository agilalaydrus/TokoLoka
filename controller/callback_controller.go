package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/middleware"
	"main.go/service"
	"math/rand"
	"net/http"
	"time"
)

type CallbackController struct {
	service service.TransactionsService
}

func NewCallbackController(service service.TransactionsService) *CallbackController {
	return &CallbackController{service: service}
}

func (cc *CallbackController) CallbackTransactionStatus(c *gin.Context) {
	middleware.Logger.Info("Controller: CallbackTransactionStatus called")

	var callbackRequest struct {
		RequestID     uint   `json:"request_id"`
		TransactionID uint   `json:"transaction_id"`
		Status        string `json:"status"`
	}

	if err := c.ShouldBindJSON(&callbackRequest); err != nil {
		middleware.Logger.Error("Error binding callback request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback request"})
		return
	}

	// Update status transaksi berdasarkan callback
	if err := cc.service.UpdateTransactionStatus(callbackRequest.TransactionID, callbackRequest.Status); err != nil {
		middleware.Logger.Error("Failed to update transaction status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	// Buat serial number acak
	serialNumber := generateSerialNumber()

	// Ambil detail transaksi setelah update
	transaction, err := cc.service.GetTransactionByID(callbackRequest.TransactionID)
	if err != nil {
		middleware.Logger.Error("Failed to fetch transaction after callback", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	// Struktur respons sukses
	response := gin.H{
		"request_id":    callbackRequest.RequestID,
		"status":        callbackRequest.Status,
		"serial_number": serialNumber,
		"message":       "Transaction status updated successfully",
		"transaction": gin.H{
			"id":          transaction.ID,
			"user_id":     transaction.UserID,
			"total_price": transaction.TotalPrice,
			"status":      transaction.Status,
			"created_at":  transaction.CreatedAt,
			"updated_at":  transaction.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// Helper untuk membuat serial number
func generateSerialNumber() string {
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	serial := make([]byte, 10)

	for i := range serial {
		serial[i] = chars[rand.Intn(len(chars))]
	}

	return string(serial)
}
