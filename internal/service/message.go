package service

import (
	"errors"
	"strings"

	"message-push-system/internal/model"
	"message-push-system/internal/repository"
)

// MessageService 消息服务接口
type MessageService interface {
	Create(req *CreateMessageRequest) (*model.Message, error)
	UpdateStatus(id int64, status string) error
	GetByID(id int64) (*model.Message, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	groupRepo   repository.GroupRepository
}

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	GroupID int64  `json:"group_id" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// NewMessageService 创建消息服务实例
func NewMessageService(messageRepo repository.MessageRepository, groupRepo repository.GroupRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
	}
}

// Create 创建消息
func (s *messageService) Create(req *CreateMessageRequest) (*model.Message, error) {
	// 验证必填字段
	if req.GroupID == 0 {
		return nil, errors.New("group_id is required")
	}
	if strings.TrimSpace(req.Subject) == "" {
		return nil, errors.New("subject cannot be empty")
	}
	if strings.TrimSpace(req.Body) == "" {
		return nil, errors.New("body cannot be empty")
	}

	// 检查群组是否存在
	_, err := s.groupRepo.GetByID(req.GroupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	message := &model.Message{
		GroupID: req.GroupID,
		Subject: req.Subject,
		Body:    req.Body,
		Status:  model.StatusPending,
	}

	return s.messageRepo.Create(message)
}

// UpdateStatus 更新消息状态
func (s *messageService) UpdateStatus(id int64, status string) error {
	return s.messageRepo.UpdateStatus(id, status)
}

// GetByID 根据 ID 获取消息
func (s *messageService) GetByID(id int64) (*model.Message, error) {
	return s.messageRepo.GetByID(id)
}
