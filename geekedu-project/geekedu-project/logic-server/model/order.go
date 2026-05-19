package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo   string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"orderNo"`
	UserID    int64      `gorm:"index;not null" json:"userId"`
	CourseID  int64      `gorm:"index;not null" json:"courseId"`
	Amount    float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayType   int        `gorm:"type:tinyint" json:"payType"`                // 1-余额 2-支付宝 3-微信
	Status    int        `gorm:"type:tinyint;default:0;index" json:"status"` // 0-待支付 1-已支付 2-已取消 3-已退款
	PayTime   *time.Time `json:"payTime"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (Order) TableName() string {
	return "order"
}
