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

	middleware.Logger.Info("User successfully registered", zap.String("email", newUser.Email))
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
		middleware.Logger.Error("Login failed", zap.String("email", loginData.Email), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Login successful", zap.String("email", loginData.Email))
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

	// Panggil service untuk mendapatkan detail user
	user, err := uc.userService.GetUserByID(userID.(uint))
	if err != nil {
		middleware.Logger.Error("Error retrieving user", zap.Uint("user_id", userID.(uint)), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	middleware.Logger.Info("User details retrieved successfully", zap.Uint("user_id", userID.(uint)))
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// UpdateUser untuk memperbarui data pengguna
func (uc *UserController) UpdateUser(c *gin.Context) {
	middleware.Logger.Info("Controller: UpdateUser called")

	userID, exists := c.Get("user_id")
	if !exists {
		middleware.Logger.Warn("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var updatedData service.UserUpdateRequest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		middleware.Logger.Error("Error binding updated user data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Update user data via service
	if err := uc.userService.UpdateUser(userID.(uint), updatedData); err != nil {
		middleware.Logger.Error("Error updating user", zap.Uint("user_id", userID.(uint)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("User information updated successfully", zap.Uint("user_id", userID.(uint)))
	c.JSON(http.StatusOK, gin.H{"message": "User information updated"})
}
