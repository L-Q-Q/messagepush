package service

import (
	"message-push-system/internal/model"
	"message-push-system/internal/repository"
)

// LogService 推送日志服务接口
type LogService interface {
	Create(log *model.PushLog) error
	ListByMessage(messageID int64) ([]*model.PushLog, error)
	List() ([]*model.PushLog, error)
}

type logService struct {
	logRepo repository.LogRepository
}

// NewLogService 创建推送日志服务实例
func NewLogService(logRepo repository.LogRepository) LogService {
	return &logService{logRepo: logRepo}
}

// Create 创建推送日志
func (s *logService) Create(log *model.PushLog) error {
	return s.logRepo.Create(log)
}

// ListByMessage 根据消息 ID 列出推送日志
func (s *logService) ListByMessage(messageID int64) ([]*model.PushLog, error) {
	return s.logRepo.ListByMessage(messageID)
}

// List 列出所有推送日志
func (s *logService) List() ([]*model.PushLog, error) {
	return s.logRepo.List()
}
