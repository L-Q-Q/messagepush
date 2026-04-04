package repository

import (
	"message-push-system/internal/model"

	"gorm.io/gorm"
)

// LogRepository 推送日志仓储接口
type LogRepository interface {
	Create(log *model.PushLog) error
	ListByMessage(messageID int64) ([]*model.PushLog, error)
	List() ([]*model.PushLog, error)
}

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository 创建推送日志仓储实例
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// Create 创建推送日志
func (r *logRepository) Create(log *model.PushLog) error {
	return r.db.Create(log).Error
}

// ListByMessage 根据消息 ID 列出推送日志
func (r *logRepository) ListByMessage(messageID int64) ([]*model.PushLog, error) {
	var logs []*model.PushLog
	if err := r.db.Where("message_id = ?", messageID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// List 列出所有推送日志
func (r *logRepository) List() ([]*model.PushLog, error) {
	var logs []*model.PushLog
	if err := r.db.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
