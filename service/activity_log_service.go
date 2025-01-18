package service

import (
	"go.uber.org/zap"
	"main.go/entity"
	"main.go/middleware"
	"main.go/repository"
	"time"
)

type ActivityLogService interface {
	CreateActivityLog(userID uint, action string, details string) error
}

type activityLogService struct {
	repository repository.ActivityLogRepository
}

// NewActivityLogService - Membuat instance baru ActivityLogService
func NewActivityLogService(repo repository.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repository: repo}
}

func (s *activityLogService) CreateActivityLog(userID uint, action string, details string) error {
	log := entity.ActivityLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return s.repository.Create(&log)
}

func (s *transactionsService) logActivity(userID uint, action string, details string) {
	if err := s.activityLogService.CreateActivityLog(userID, action, details); err != nil {
		middleware.Logger.Error("Failed to create activity log", zap.Error(err))
	}
}
