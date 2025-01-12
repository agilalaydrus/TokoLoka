package controller

import (
	"github.com/gin-gonic/gin"
	"log"
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
	var newUser service.UserRegisterRequest

	// Bind JSON request body ke dalam struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Printf("Error while binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Panggil service untuk registrasi user
	if err := uc.userService.RegisterUser(newUser); err != nil {
		log.Printf("Error while registering user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User successfully registered"})
}

// LoginUser untuk melakukan login dan menghasilkan JWT token
func (uc *UserController) LoginUser(c *gin.Context) {
	var loginData service.UserLoginRequest

	// Bind JSON request body ke dalam struct
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Printf("Error while binding login data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Panggil service untuk login dan verifikasi pengguna
	token, err := uc.userService.LoginUser(loginData)
	if err != nil {
		log.Printf("Login failed for email %s: %v", loginData.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Kirim token ke client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// GetUserDetails untuk mengambil detail pengguna berdasarkan ID
func (uc *UserController) GetUserDetails(c *gin.Context) {
	// Ambil user_id dari JWT yang sudah didecode sebelumnya
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Panggil service untuk mendapatkan detail user
	user, err := uc.userService.GetUserByID(userID.(uint))
	if err != nil {
		log.Printf("Error retrieving user %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Kirimkan detail pengguna
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var updatedData service.UserUpdateRequest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		log.Printf("Error binding updated user data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Update user data via service
	if err := uc.userService.UpdateUser(userID.(uint), updatedData); err != nil {
		log.Printf("Error updating user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User information updated"})
}
