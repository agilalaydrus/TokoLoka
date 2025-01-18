package service

import (
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
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (s *UserService) RegisterUser(user UserRegisterRequest) error {
	middleware.Logger.Info("Service: Registering user", zap.String("email", user.Email))

	// Validasi input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		middleware.Logger.Warn("Service: Missing required fields", zap.Any("user", user))
		return middleware.NewAppError(400, "Username, email, and password are required", nil)
	}

	if !isValidEmail(user.Email) {
		middleware.Logger.Warn("Service: Invalid email format", zap.String("email", user.Email))
		return middleware.NewAppError(400, "Invalid email format", nil)
	}

	// Default role jika kosong
	if user.Role == "" {
		user.Role = "user"
	}

	// Cek email sudah digunakan
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		middleware.Logger.Warn("Service: Email already in use", zap.String("email", user.Email))
		return middleware.NewAppError(409, "Email already in use", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.Logger.Error("Service: Failed to hash password", zap.Error(err))
		return middleware.NewAppError(500, "Failed to hash password", err)
	}

	newUser := entity.User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     user.Role,
	}

	if err := s.userRepo.Create(&newUser); err != nil {
		middleware.Logger.Error("Service: Failed to create user", zap.Error(err))
		return middleware.NewAppError(500, "Failed to create user", err)
	}

	middleware.Logger.Info("Service: User registered successfully", zap.Uint("user_id", newUser.ID))
	return nil
}

func (s *UserService) LoginUser(user UserLoginRequest) (string, error) {
	middleware.Logger.Info("Service: Logging in user", zap.String("email", user.Email))

	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		middleware.Logger.Warn("Service: Invalid credentials", zap.String("email", user.Email))
		return "", middleware.NewAppError(401, "Invalid credentials", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		middleware.Logger.Warn("Service: Invalid password", zap.String("email", user.Email))
		return "", middleware.NewAppError(401, "Invalid credentials", nil)
	}

	token, err := generateJWT(existingUser)
	if err != nil {
		middleware.Logger.Error("Service: Failed to generate JWT", zap.Error(err))
		return "", middleware.NewAppError(500, "Failed to generate token", err)
	}

	middleware.Logger.Info("Service: User logged in successfully", zap.Uint("user_id", existingUser.ID))
	return token, nil
}

func generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		middleware.Logger.Fatal("Service: JWT_SECRET is not set")
		return "", middleware.NewAppError(500, "JWT_SECRET is not set", nil)
	}

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		middleware.Logger.Error("Service: Failed to sign JWT", zap.Error(err))
		return "", middleware.NewAppError(500, "Failed to sign JWT", err)
	}

	return signedToken, nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (s *UserService) UpdateUser(userID uint, updatedData UserUpdateRequest) error {
	middleware.Logger.Info("Service: Updating user", zap.Uint("user_id", userID))

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		middleware.Logger.Warn("Service: User not found", zap.Uint("user_id", userID))
		return middleware.NewAppError(404, "User not found", err)
	}

	if updatedData.Username != "" {
		user.Username = updatedData.Username
	}
	if updatedData.Email != "" {
		existingUser, _ := s.userRepo.FindByEmail(updatedData.Email)
		if existingUser != nil && existingUser.ID != userID {
			middleware.Logger.Warn("Service: Email already in use", zap.String("email", updatedData.Email))
			return middleware.NewAppError(409, "Email already in use", nil)
		}
		user.Email = updatedData.Email
	}
	if updatedData.Role != "" {
		user.Role = updatedData.Role
	}

	if err := s.userRepo.Update(user); err != nil {
		middleware.Logger.Error("Service: Failed to update user", zap.Error(err))
		return middleware.NewAppError(500, "Failed to update user", err)
	}

	middleware.Logger.Info("Service: User updated successfully", zap.Uint("user_id", userID))
	return nil
}

func (s *UserService) GetUserByID(userID uint) (*entity.User, error) {
	middleware.Logger.Info("Service: Fetching user by ID", zap.Uint("user_id", userID))

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		middleware.Logger.Warn("Service: User not found", zap.Uint("user_id", userID))
		return nil, middleware.NewAppError(404, "User not found", err)
	}

	middleware.Logger.Info("Service: User fetched successfully", zap.Any("user", user))
	return user, nil
}
