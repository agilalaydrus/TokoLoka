package repository

import (
	"gorm.io/gorm"
	"main.go/entity"
)

type TokenRepository interface {
	Save(token *entity.RefreshToken) error
	FindByToken(token string) (*entity.RefreshToken, error)
	Delete(token string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) Save(token *entity.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) FindByToken(token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	return &rt, err
}

func (r *tokenRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}
