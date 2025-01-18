package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Helper function untuk membuat AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ErrorHandler middleware untuk menangani error secara terpusat
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Jalankan handler berikutnya
		c.Next()

		// Tangkap error jika ada
		if len(c.Errors) > 0 {
			// Ambil error terakhir
			err := c.Errors.Last().Err

			// Jika error adalah AppError, gunakan informasinya
			if appErr, ok := err.(*AppError); ok {
				Logger.Error("AppError occurred", zap.Error(appErr.Err), zap.String("message", appErr.Message))
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
				return
			}

			// Jika error biasa, log dan kirimkan response generik
			Logger.Error("Unhandled error occurred", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
