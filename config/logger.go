package config

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger initializes the Zap logger
func InitLogger() error {
	var err error
	// Gunakan konfigurasi bawaan Zap (Production Mode)
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(Logger) // Menjadikan logger ini sebagai global
	return nil
}
