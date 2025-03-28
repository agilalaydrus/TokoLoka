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

// ======================== REGISTER ==========================
func (uc *UserController) RegisterUser(c *gin.Context) {
	var newUser service.UserRegisterRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		middleware.Logger.Error("Error binding request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := uc.userService.RegisterUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("User registered", zap.String("phone", newUser.PhoneNumber))
	c.JSON(http.StatusCreated, gin.H{"message": "User successfully registered"})
}

// ======================== LOGIN ==========================
func (uc *UserController) LoginUser(c *gin.Context) {
	var loginData service.UserLoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		middleware.Logger.Error("Error binding login data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := uc.userService.LoginUser(loginData)
	if err != nil {
		middleware.Logger.Warn("Login failed", zap.String("phone", loginData.PhoneNumber), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Info("Login successful", zap.String("phone", loginData.PhoneNumber))
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ======================== REFRESH TOKEN ==========================
func (uc *UserController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	newAccessToken, err := uc.userService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// ======================== GET PROFILE ==========================
func (uc *UserController) GetUserDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.userService.GetUserByID(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"full_name":  user.FullName,
		"phone":      user.PhoneNumber,
		"email":      user.Email,
		"address":    user.Address,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// ======================== UPDATE PROFILE ==========================
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedData service.UserUpdateRequest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := uc.userService.UpdateUser(userIDUint, updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// ======================== UPDATE PROFILE ==========================
func (uc *UserController) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	err := uc.userService.LogoutUser(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
