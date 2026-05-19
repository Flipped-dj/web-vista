package database

import (
	"fmt"
	"log"

	"geekedu/logic-server/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(host, user, password, dbname string, port int) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Course{},
		&model.Video{},
		&model.Order{},
		&model.Chapter{},
		&model.Category{},
		&model.UserCourse{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}
