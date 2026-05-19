package model

import (
	"time"
)

// User 用户模型（精简版，只保留必要字段）
type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname  string     `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string     `gorm:"type:varchar(500)" json:"avatar"`
	Role      int        `gorm:"type:tinyint;default:0" json:"role"`   // 0-学员 1-讲师 2-管理员
	Status    int        `gorm:"type:tinyint;default:1" json:"status"` // 0-禁用 1-正常
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
