package repository

import (
	"message-push-system/internal/model"

	"gorm.io/gorm"
)

// MemberRepository 成员仓储接口
type MemberRepository interface {
	Create(member *model.Member) (*model.Member, error)
	Delete(id int64) error
	ListByGroup(groupID int64) ([]*model.Member, error)
	GetByID(id int64) (*model.Member, error)
}

type memberRepository struct {
	db *gorm.DB
}

// NewMemberRepository 创建成员仓储实例
func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

// Create 创建成员
func (r *memberRepository) Create(member *model.Member) (*model.Member, error) {
	if err := r.db.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

// Delete 删除成员
func (r *memberRepository) Delete(id int64) error {
	return r.db.Delete(&model.Member{}, id).Error
}

// ListByGroup 列出群组的所有成员
func (r *memberRepository) ListByGroup(groupID int64) ([]*model.Member, error) {
	var members []*model.Member
	if err := r.db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// GetByID 根据 ID 获取成员
func (r *memberRepository) GetByID(id int64) (*model.Member, error) {
	var member model.Member
	if err := r.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}
