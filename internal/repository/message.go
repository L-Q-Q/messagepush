package repository

import (
	"message-push-system/internal/model"

	"gorm.io/gorm"
)

// MessageRepository 消息仓储接口
type MessageRepository interface {
	Create(message *model.Message) (*model.Message, error)
	UpdateStatus(id int64, status string) error
	GetByID(id int64) (*model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓储实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Create 创建消息
func (r *messageRepository) Create(message *model.Message) (*model.Message, error) {
	if err := r.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// UpdateStatus 更新消息状态
func (r *messageRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Message{}).Where("id = ?", id).Update("status", status).Error
}

// GetByID 根据 ID 获取消息
func (r *messageRepository) GetByID(id int64) (*model.Message, error) {
	var message model.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}
