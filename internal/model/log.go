package model

import "time"

// PushLog 推送日志模型
type PushLog struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	MessageID    int64     `gorm:"not null;index" json:"message_id"`
	Recipient    string    `gorm:"type:varchar(255);not null" json:"recipient"`
	Status       string    `gorm:"type:varchar(50);not null;index" json:"status"`
	ErrorMessage string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (PushLog) TableName() string {
	return "push_logs"
}
