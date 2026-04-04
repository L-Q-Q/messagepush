package service

import (
	"errors"
	"strings"

	"message-push-system/internal/model"
	"message-push-system/internal/repository"
)

// GroupService 群组服务接口
type GroupService interface {
	Create(req *CreateGroupRequest) (*model.Group, error)
	Delete(id int64) error
	List() ([]*model.Group, error)
	GetByID(id int64) (*model.Group, error)
}

type groupService struct {
	groupRepo repository.GroupRepository
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Name         string `json:"name" binding:"required"`
	SMTPServer   string `json:"smtp_server" binding:"required"`
	SMTPPort     int    `json:"smtp_port" binding:"required"`
	SMTPUsername string `json:"smtp_username" binding:"required"`
	SMTPPassword string `json:"smtp_password" binding:"required"`
}

// NewGroupService 创建群组服务实例
func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &groupService{groupRepo: groupRepo}
}

// Create 创建群组
func (s *groupService) Create(req *CreateGroupRequest) (*model.Group, error) {
	// 验证群组名称
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("group name cannot be empty")
	}

	// 验证 SMTP 配置
	if strings.TrimSpace(req.SMTPServer) == "" {
		return nil, errors.New("smtp server cannot be empty")
	}
	if req.SMTPPort < 1 || req.SMTPPort > 65535 {
		return nil, errors.New("smtp port must be between 1 and 65535")
	}
	if strings.TrimSpace(req.SMTPUsername) == "" {
		return nil, errors.New("smtp username cannot be empty")
	}
	if strings.TrimSpace(req.SMTPPassword) == "" {
		return nil, errors.New("smtp password cannot be empty")
	}

	group := &model.Group{
		Name:         req.Name,
		SMTPServer:   req.SMTPServer,
		SMTPPort:     req.SMTPPort,
		SMTPUsername: req.SMTPUsername,
		SMTPPassword: req.SMTPPassword,
	}

	return s.groupRepo.Create(group)
}

// Delete 删除群组
func (s *groupService) Delete(id int64) error {
	return s.groupRepo.Delete(id)
}

// List 列出所有群组
func (s *groupService) List() ([]*model.Group, error) {
	return s.groupRepo.List()
}

// GetByID 根据 ID 获取群组
func (s *groupService) GetByID(id int64) (*model.Group, error) {
	return s.groupRepo.GetByID(id)
}
