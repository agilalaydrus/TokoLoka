package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main.go/middleware"
	"main.go/service"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// RegisterUser untuk melakukan registrasi pengguna baru
func (uc *UserController) RegisterUser(c *gin.Context) {
	middleware.Logger.Info("Controller: RegisterUser called")

	var newUser service.UserRegisterRequest

	// Bind JSON request body ke dalam struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		middleware.Logger.Error("Error binding request data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Panggil service untuk registrasi user
	if err := uc.userService.RegisterUser(newUser); err != nil {
		middleware.Logger.Error("Error registering user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("User successfully registered", zap.String("phone_number", newUser.PhoneNumber))
	c.JSON(http.StatusCreated, gin.H{"message": "User successfully registered"})
}

// LoginUser untuk melakukan login dan menghasilkan JWT token
func (uc *UserController) LoginUser(c *gin.Context) {
	middleware.Logger.Info("Controller: LoginUser called")

	var loginData service.UserLoginRequest

	// Bind JSON request body ke dalam struct
	if err := c.ShouldBindJSON(&loginData); err != nil {
		middleware.Logger.Error("Error binding login data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Panggil service untuk login dan verifikasi pengguna
	token, err := uc.userService.LoginUser(loginData)
	if err != nil {
		middleware.Logger.Error("Login failed", zap.String("phone_number", loginData.PhoneNumber), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Login successful", zap.String("phone_number", loginData.PhoneNumber))
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// GetUserDetails untuk mengambil detail pengguna berdasarkan ID
func (uc *UserController) GetUserDetails(c *gin.Context) {
	middleware.Logger.Info("Controller: GetUserDetails called")

	// Ambil user_id dari JWT yang sudah didecode sebelumnya
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.Logger.Warn("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Safe casting to uint
	userIDUint, ok := userID.(uint)
	if !ok {
		middleware.Logger.Error("Failed to cast user_id to uint")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Panggil service untuk mendapatkan detail user
	user, err := uc.userService.GetUserByID(userIDUint)
	if err != nil {
		middleware.Logger.Error("Error retrieving user", zap.Uint("user_id", userIDUint), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	middleware.Logger.Info("User details retrieved successfully", zap.Uint("user_id", userIDUint))
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"full_name":  user.FullName,    // Menggunakan full_name sebagai pengganti username
		"phone":      user.PhoneNumber, // Nomor HP sebagai identifier utama
		"email":      user.Email,       // Opsional, jika masih digunakan
		"address":    user.Address,     // Menambahkan alamat lengkap
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// UpdateUser untuk memperbarui data pengguna
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var updatedData service.UserUpdateRequest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := uc.userService.UpdateUser(userIDUint, updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
