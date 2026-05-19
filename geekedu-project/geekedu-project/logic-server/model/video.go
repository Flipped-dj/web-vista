package model

import (
	"time"
)

// Video 视频模型
type Video struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  int64      `gorm:"index;not null" json:"courseId"`
	ChapterID *int64     `gorm:"index" json:"chapterId"`
	Title     string     `gorm:"type:varchar(200);not null" json:"title"`
	Duration  int        `gorm:"default:0" json:"duration"`
	VideoURL  string     `gorm:"type:varchar(500)" json:"videoUrl"`
	Cover     string     `gorm:"type:varchar(255)" json:"cover"`
	IsFree    int        `gorm:"type:tinyint;default:0" json:"isFree"`
	Sort      int        `gorm:"default:0" json:"sort"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (Video) TableName() string {
	return "video"
}

// Chapter 章节模型
type Chapter struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  int64      `gorm:"index;not null" json:"courseId"`
	Title     string     `gorm:"type:varchar(200);not null" json:"title"`
	Sort      int        `gorm:"default:0" json:"sort"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (Chapter) TableName() string {
	return "chapter"
}

// Category 分类模型
type Category struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:varchar(50);not null" json:"name"`
	ParentID  int64      `gorm:"index;default:0" json:"parentId"`
	Sort      int        `gorm:"default:0" json:"sort"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "category"
}

// UserCourse 用户课程关系表
type UserCourse struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"userId"`
	CourseID  int64     `gorm:"index;not null" json:"courseId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (UserCourse) TableName() string {
	return "user_course"
}
