package model

import "time"

// Member 成员模型
type Member struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	GroupID   int64     `gorm:"not null;index" json:"group_id"`
	Email     string    `gorm:"type:varchar(255);not null;index" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Member) TableName() string {
	return "members"
}
