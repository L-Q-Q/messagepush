package model

import "time"

// Message 消息模型
type Message struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	GroupID   int64     `gorm:"not null;index" json:"group_id"`
	Subject   string    `gorm:"type:varchar(500);not null" json:"subject"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Status    string    `gorm:"type:varchar(50);not null;default:pending;index" json:"status"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// 消息状态常量
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
)
