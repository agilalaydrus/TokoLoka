package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var JWT_SECRET []byte

// Init fungsi untuk memuat JWT_SECRET
func init() {
	// Muat file .env jika tersedia
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables.")
	}

	// Ambil nilai JWT_SECRET dari environment variable
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	if len(JWT_SECRET) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set")
	} else {
		log.Println("JWT_SECRET berhasil dimuat")
	}
}

// AuthorizeJWT middleware untuk memverifikasi JWT dan menyimpan klaim ke context
func AuthorizeJWT(c *gin.Context) {
	log.Println("Middleware: AuthorizeJWT called")

	// Ambil token dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Middleware: Header Authorization tidak ditemukan")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization tidak ditemukan"})
		c.Abort()
		return
	}

	// Periksa format token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Middleware: Format token tidak valid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
		c.Abort()
		return
	}

	// Ambil token setelah kata "Bearer "
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse dan verifikasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi metode tanda tangan
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Middleware: Metode tanda tangan tidak valid")
			return nil, fmt.Errorf("metode tanda tangan tidak valid")
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		log.Printf("Middleware: Error parsing token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}

	if !token.Valid {
		log.Println("Middleware: Token tidak valid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}

	// Ambil klaim dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Middleware: Token claims tidak valid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims tidak valid"})
		c.Abort()
		return
	}

	// Ambil user_id dari klaim (menggunakan sub)
	userID, err := extractUserID(claims["sub"])
	if err != nil {
		log.Println("Middleware: Error extracting user_id from token", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id tidak valid dalam token"})
		c.Abort()
		return
	}

	// Ambil role dari klaim
	role, roleOk := claims["role"].(string)
	if !roleOk || role == "" {
		log.Println("Middleware: Role tidak ditemukan dalam token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan dalam token"})
		c.Abort()
		return
	}

	// Simpan user_id dan role ke context untuk digunakan di handler berikutnya
	c.Set("user_id", userID)
	c.Set("role", role)

	log.Printf("Middleware: Token berhasil diverifikasi (user_id=%v, role=%s)", userID, role)

	// Lanjutkan ke request berikutnya
	c.Next()
}

// Fungsi untuk mengekstrak user_id dari klaim
func extractUserID(sub interface{}) (uint, error) {
	switch v := sub.(type) {
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("error converting string to uint")
		}
		return uint(id), nil
	default:
		return 0, fmt.Errorf("invalid type for user_id")
	}
}

func RoleBasedAccessControl(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context
		role, exists := c.Get("role")
		if !exists {
			log.Println("Middleware: Role tidak ditemukan di context")
			c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak ditemukan"})
			c.Abort()
			return
		}

		// Periksa apakah user memiliki akses
		currentRole := role.(string)
		if currentRole != requiredRole && requiredRole != "ANY" {
			log.Printf("Middleware: Access denied (required=%s, current=%s)", requiredRole, currentRole)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for your role"})
			c.Abort()
			return
		}

		log.Printf("Middleware: Access granted for role=%s", currentRole)

		// Lanjutkan jika role sesuai
		c.Next()
	}
}

// HasRole - Helper function untuk memvalidasi role user
func HasRole(c *gin.Context, roles []string) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	for _, r := range roles {
		if role == r {
			return true
		}
	}

	return false
}
