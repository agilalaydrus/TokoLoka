package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.Logger

func InitLogger() {
	// Konfigurasi Zap
	logFile, _ := os.Create("logs/app.log") // Buat file log
	writeSyncer := zapcore.AddSync(logFile)

	// Konfigurasi Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Konfigurasi Core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                // Format log JSON
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer), // Tulis ke file dan konsol
		zapcore.DebugLevel, // Level log minimum
	)

	// Buat logger
	Logger = zap.New(core, zap.AddCaller())
	Logger.Info("Zap logger initialized with file logging")
}

// RequestLogger middleware untuk mencatat setiap permintaan HTTP
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now() // Waktu mulai

		// Proses request
		c.Next()

		// Data log
		statusCode := c.Writer.Status()
		latency := time.Since(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Logging berdasarkan status
		if statusCode >= 500 {
			Logger.Error("Internal Server Error",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage))
		} else if statusCode >= 400 {
			Logger.Warn("Client Error",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage))
		} else {
			Logger.Info("Request",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Duration("latency", latency))
		}
	}
}
