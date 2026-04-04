package service

import (
	"errors"
	"regexp"
	"strings"

	"message-push-system/internal/model"
	"message-push-system/internal/repository"
)

// MemberService 成员服务接口
type MemberService interface {
	Add(groupID int64, email string) (*model.Member, error)
	Remove(id int64) error
	ListByGroup(groupID int64) ([]*model.Member, error)
}

type memberService struct {
	memberRepo repository.MemberRepository
}

// NewMemberService 创建成员服务实例
func NewMemberService(memberRepo repository.MemberRepository) MemberService {
	return &memberService{memberRepo: memberRepo}
}

// Add 添加成员
func (s *memberService) Add(groupID int64, email string) (*model.Member, error) {
	// 验证邮箱
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	// 简单的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	member := &model.Member{
		GroupID: groupID,
		Email:   email,
	}

	return s.memberRepo.Create(member)
}

// Remove 删除成员
func (s *memberService) Remove(id int64) error {
	return s.memberRepo.Delete(id)
}

// ListByGroup 列出群组的所有成员
func (s *memberService) ListByGroup(groupID int64) ([]*model.Member, error) {
	return s.memberRepo.ListByGroup(groupID)
}
