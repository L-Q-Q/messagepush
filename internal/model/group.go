package model

import "time"

// Group 群组模型
type Group struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	SMTPServer   string    `gorm:"type:varchar(255);not null" json:"smtp_server"`
	SMTPPort     int       `gorm:"not null" json:"smtp_port"`
	SMTPUsername string    `gorm:"type:varchar(255);not null" json:"smtp_username"`
	SMTPPassword string    `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}
