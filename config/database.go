package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"main.go/entity"
	"os"
)

// Global DB instance
var DB *gorm.DB

// LoadEnv menginisialisasi file .env untuk mengambil variabel environment
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}
	return nil
}

// InitDB menginisialisasi koneksi ke database MySQL dan melakukan migrasi tabel
func InitDB() error {
	// Memuat konfigurasi dari .env
	if err := LoadEnv(); err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}

	// Ambil konfigurasi dari environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Membuat string koneksi database
	databaseURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Menginisialisasi koneksi dengan GORM
	var err error
	DB, err = gorm.Open(mysql.Open(databaseURI), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %w", err)
	}

	// Mengecek koneksi ke database
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("Gagal mendapatkan database connection: %w", err)
	}

	// Menjalankan ping untuk memastikan koneksi berhasil
	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("Gagal melakukan ping ke database: %w", err)
	}

	log.Println("Koneksi ke database berhasil!")

	// Melakukan migrasi tabel
	err = DB.AutoMigrate(
		&entity.User{},
		&entity.Product{},
		&entity.Transaction{},
		&entity.TransactionItem{},
		&entity.Category{},
	)
	if err != nil {
		return fmt.Errorf("Gagal melakukan migrasi: %w", err)
	}

	log.Println("Migrasi tabel berhasil!")
	return nil
}
