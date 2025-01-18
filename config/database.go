package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"main.go/entity"
)

var DB *gorm.DB

func InitDB() error {
	// Ambil konfigurasi dari environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Membuat string koneksi database
	databaseURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Menginisialisasi koneksi dengan GORM
	var err error
	DB, err = gorm.Open(mysql.Open(databaseURI), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke database: %w", err)
	}

	log.Println("Koneksi ke database berhasil!")

	// Jalankan migration untuk semua tabel
	err = DB.AutoMigrate(
		&entity.User{},
		&entity.Product{},
		&entity.Category{},
		&entity.Transaction{},
		&entity.TransactionItem{},
		&entity.ActivityLog{},
		&entity.ReportLog{},
	)
	if err != nil {
		return fmt.Errorf("gagal melakukan migrasi: %w", err)
	}

	log.Println("Migrasi tabel berhasil!")
	return nil
}
