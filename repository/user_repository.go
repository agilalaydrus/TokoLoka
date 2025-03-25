package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main.go/entity"
	"main.go/middleware"
)

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	FindByPhoneNumber(phone string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	middleware.Logger.Info("Repository: Fetching user by email", zap.String("email", email))
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		middleware.Logger.Warn("Repository: User not found", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	middleware.Logger.Info("Repository: User fetched successfully", zap.Any("user", user))
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	middleware.Logger.Info("Repository: Fetching user by ID", zap.Uint("user_id", id))
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		middleware.Logger.Warn("Repository: User not found", zap.Uint("user_id", id), zap.Error(err))
		return nil, err
	}
	middleware.Logger.Info("Repository: User fetched successfully", zap.Any("user", user))
	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	middleware.Logger.Info("Repository: Creating user", zap.Any("user", user))
	if err := r.db.Create(user).Error; err != nil {
		middleware.Logger.Error("Repository: Error creating user", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: User created successfully", zap.Uint("user_id", user.ID))
	return nil
}

func (r *userRepository) Update(user *entity.User) error {
	middleware.Logger.Info("Repository: Updating user", zap.Uint("user_id", user.ID))
	if err := r.db.Save(user).Error; err != nil {
		middleware.Logger.Error("Repository: Error updating user", zap.Error(err))
		return err
	}
	middleware.Logger.Info("Repository: User updated successfully", zap.Uint("user_id", user.ID))
	return nil
}

func (r *userRepository) FindByPhoneNumber(phone string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
