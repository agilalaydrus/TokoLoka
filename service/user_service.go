package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"main.go/entity"
	"main.go/repository"
	"os"
	"regexp"
	"time"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Struktur request untuk registrasi pengguna
type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // Role akan diberikan default jika kosong
}

// Struktur request untuk login pengguna
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struktur untuk memperbarui data pengguna
type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // Optional jika pengguna dapat mengubah role
}

// Fungsi untuk registrasi pengguna
func (s *UserService) RegisterUser(user UserRegisterRequest) error {
	// Cek apakah email sudah digunakan
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already in use")
	}

	// Validasi input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("username, email, and password are required")
	}

	// Validasi format email
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Defaultkan role jika tidak diisi
	if user.Role == "" {
		user.Role = "user" // Role default
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Simpan pengguna baru
	newUser := entity.User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     user.Role,
	}

	if err := s.userRepo.Create(&newUser); err != nil {
		return err
	}

	return nil
}

// Fungsi untuk login pengguna
func (s *UserService) LoginUser(user UserLoginRequest) (string, error) {
	// Cari pengguna berdasarkan email
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWT(existingUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Fungsi untuk menghasilkan token JWT
func generateJWT(user *entity.User) (string, error) {
	// Set JWT claims
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Pastikan JWT_SECRET ada di environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	// Sign token
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Fungsi untuk validasi format email
func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (s *UserService) UpdateUser(userID uint, updatedData UserUpdateRequest) error {
	// Cari pengguna berdasarkan ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update hanya kolom yang diisi
	if updatedData.Username != "" {
		user.Username = updatedData.Username
	}
	if updatedData.Email != "" {
		// Cek apakah email sudah digunakan oleh pengguna lain
		existingUser, _ := s.userRepo.FindByEmail(updatedData.Email)
		if existingUser != nil && existingUser.ID != userID {
			return errors.New("email already in use")
		}
		user.Email = updatedData.Email
	}
	if updatedData.Role != "" {
		user.Role = updatedData.Role
	}

	// Simpan perubahan
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

func (s *UserService) GetUserByID(userID uint) (*entity.User, error) {
	// Cari pengguna berdasarkan ID melalui repository
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
