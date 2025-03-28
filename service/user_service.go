package service

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/entity"
	"main.go/middleware"
	"main.go/repository"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

type UserRegisterRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address"`
}

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserUpdateRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

// ======================== REGISTER ==========================
func (s *UserService) RegisterUser(user UserRegisterRequest) error {
	if user.FullName == "" || user.PhoneNumber == "" || user.Password == "" {
		return middleware.NewAppError(400, "Full name, phone number, and password are required", nil)
	}

	if user.Email != "" && !isValidEmail(user.Email) {
		return middleware.NewAppError(400, "Invalid email format", nil)
	}

	if existing, _ := s.userRepo.FindByPhoneNumber(user.PhoneNumber); existing != nil {
		return middleware.NewAppError(409, "Phone number already registered", nil)
	}

	if user.Email != "" {
		if existing, _ := s.userRepo.FindByEmail(user.Email); existing != nil {
			return middleware.NewAppError(409, "Email already registered", nil)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return middleware.NewAppError(500, "Internal error while securing your account", err)
	}

	newUser := entity.User{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    string(hashedPassword),
		Address:     user.Address,
		Role:        "user",
	}

	if err := s.userRepo.Create(&newUser); err != nil {
		return middleware.NewAppError(500, "Failed to create user", err)
	}

	middleware.Logger.Info("Service: User registered", zap.Uint("user_id", newUser.ID))
	return nil
}

// ======================== LOGIN ==========================
func (s *UserService) LoginUser(user UserLoginRequest) (string, string, error) {
	if user.PhoneNumber == "" || user.Password == "" {
		return "", "", middleware.NewAppError(400, "Phone number and password are required", nil)
	}

	existingUser, err := s.userRepo.FindByPhoneNumber(user.PhoneNumber)
	if err != nil || existingUser == nil {
		return "", "", middleware.NewAppError(401, "Invalid phone number or password", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return "", "", middleware.NewAppError(401, "Invalid phone number or password", nil)
	}

	accessToken, err := generateJWT(existingUser)
	if err != nil {
		return "", "", middleware.NewAppError(500, "Failed to generate token", err)
	}

	refreshToken := uuid.New().String()
	refreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)

	rt := entity.RefreshToken{
		UserID:    existingUser.ID,
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpires,
	}

	if err := s.tokenRepo.Save(&rt); err != nil {
		return "", "", middleware.NewAppError(500, "Failed to store refresh token", err)
	}

	middleware.Logger.Info("Service: Login successful", zap.Uint("user_id", existingUser.ID))
	return accessToken, refreshToken, nil
}

// ======================== REFRESH TOKEN ==========================
func (s *UserService) RefreshAccessToken(refreshToken string) (string, error) {
	rt, err := s.tokenRepo.FindByToken(refreshToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return "", middleware.NewAppError(401, "Invalid or expired refresh token", err)
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return "", middleware.NewAppError(404, "User not found", err)
	}

	newAccessToken, err := generateJWT(user)
	if err != nil {
		return "", middleware.NewAppError(500, "Failed to generate token", err)
	}

	return newAccessToken, nil
}

// ======================== UPDATE ==========================
func (s *UserService) UpdateUser(userID uint, updatedData UserUpdateRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return middleware.NewAppError(404, "User not found", err)
	}

	if updatedData.FullName != "" {
		user.FullName = updatedData.FullName
	}
	if updatedData.Address != "" {
		user.Address = updatedData.Address
	}
	if updatedData.Email != "" {
		if !isValidEmail(updatedData.Email) {
			return middleware.NewAppError(400, "Invalid email format", nil)
		}
		existing, _ := s.userRepo.FindByEmail(updatedData.Email)
		if existing != nil && existing.ID != userID {
			return middleware.NewAppError(409, "Email already used by another account", nil)
		}
		user.Email = updatedData.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return middleware.NewAppError(500, "Failed to update user", err)
	}

	middleware.Logger.Info("Service: User updated successfully", zap.Uint("user_id", userID))
	return nil
}

// ======================== GET BY ID ==========================
func (s *UserService) GetUserByID(userID uint) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, middleware.NewAppError(404, "User not found", err)
	}

	return user, nil
}

// ======================== JWT & VALIDATOR ==========================
func generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", middleware.NewAppError(500, "JWT secret is not set", nil)
	}

	return token.SignedString([]byte(jwtSecret))
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ======================== LOGOUT USER ==========================
func (s *UserService) LogoutUser(refreshToken string) error {
	err := s.tokenRepo.Delete(refreshToken)
	if err != nil {
		return middleware.NewAppError(500, "Failed to revoke refresh token", err)
	}
	return nil
}
