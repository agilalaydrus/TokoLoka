package repository

import (
	"gorm.io/gorm"
	"main.go/entity"
)

type ActivityLogRepository interface {
	Create(log *entity.ActivityLog) error
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(log *entity.ActivityLog) error {
	return r.db.Create(log).Error
}
