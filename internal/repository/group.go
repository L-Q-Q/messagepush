package repository

import (
	"message-push-system/internal/model"

	"gorm.io/gorm"
)

// GroupRepository 群组仓储接口
type GroupRepository interface {
	Create(group *model.Group) (*model.Group, error)
	Delete(id int64) error
	List() ([]*model.Group, error)
	GetByID(id int64) (*model.Group, error)
}

type groupRepository struct {
	db *gorm.DB
}

// NewGroupRepository 创建群组仓储实例
func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

// Create 创建群组
func (r *groupRepository) Create(group *model.Group) (*model.Group, error) {
	if err := r.db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// Delete 删除群组（级联删除成员）
func (r *groupRepository) Delete(id int64) error {
	return r.db.Delete(&model.Group{}, id).Error
}

// List 列出所有群组
func (r *groupRepository) List() ([]*model.Group, error) {
	var groups []*model.Group
	if err := r.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetByID 根据 ID 获取群组
func (r *groupRepository) GetByID(id int64) (*model.Group, error) {
	var group model.Group
	if err := r.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
