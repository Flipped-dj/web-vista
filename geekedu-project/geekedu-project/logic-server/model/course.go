package model

import (
	"time"
)

// Course 课程模型
type Course struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string     `gorm:"type:varchar(200);not null" json:"title"`
	Cover         string     `gorm:"type:varchar(255)" json:"cover"`
	Description   string     `gorm:"type:text" json:"description"`
	CategoryID    int64      `gorm:"index" json:"categoryId"`
	TeacherID     int64      `gorm:"index" json:"teacherId"`
	Price         float64    `gorm:"type:decimal(10,2);default:0" json:"price"`
	OriginalPrice float64    `gorm:"type:decimal(10,2);default:0" json:"originalPrice"`
	StudentCount  int        `gorm:"default:0" json:"studentCount"`
	LessonCount   int        `gorm:"default:0" json:"lessonCount"`
	Duration      int        `gorm:"default:0" json:"duration"`                 // 秒
	Level         int        `gorm:"type:tinyint;default:1" json:"level"`       // 1-入门 2-初级 3-中级 4-高级
	IsRecommend   int        `gorm:"type:tinyint;default:0" json:"isRecommend"` // 0-否 1-是
	Status        int        `gorm:"type:tinyint;default:1" json:"status"`      // 0-下架 1-上架
	Sort          int        `gorm:"default:0" json:"sort"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "course"
}
