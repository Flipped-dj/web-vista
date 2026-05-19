package service

import (
	"context"
	"errors"

	"geekedu/common/config"
	"geekedu/common/oss"
	"geekedu/logic-server/database"
	"geekedu/logic-server/model"

	"gorm.io/gorm"
)

type VideoService struct {
	ossConfig *config.OSSConfig
}

func NewVideoService(ossConfig *config.OSSConfig) *VideoService {
	return &VideoService{
		ossConfig: ossConfig,
	}
}

// GetVideoPlayURL 获取视频播放地址（预签名URL）
// 核心：播放鉴权流程
func (s *VideoService) GetVideoPlayURL(ctx context.Context, videoId, userId int64) (string, int, error) {
	// 1. 查询视频信息
	var video model.Video
	err := database.DB.Where("id = ?", videoId).First(&video).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", 404, errors.New("video not found")
		}
		return "", 500, err
	}

	// 2. 查询用户角色
	var user model.User
	userErr := database.DB.Where("id = ?", userId).First(&user).Error

	// 3. 检查视频播放权限
	// 付费视频：需要购买，或者是管理员
	if video.IsFree == 0 {
		// 查询课程信息获取创建者ID
		var course model.Course
		courseErr := database.DB.Where("id = ?", video.CourseID).First(&course).Error

		// 检查是否为课程创建者（管理员可以免费观看自己的课程）
		isCreator := courseErr == nil && userErr == nil && user.Role == 2

		if !isCreator {
			// 非创建者，需要检查是否已购买
			hasPurchased, err := s.CheckUserCourse(ctx, userId, video.CourseID)
			if err != nil {
				return "", 500, err
			}
			if !hasPurchased {
				return "", 403, errors.New("you have not purchased this course")
			}
		}
	}

	// 4. 生成OSS预签名URL（有效期3600秒）
	if s.ossConfig.AccessKeyID == "" || s.ossConfig.AccessKeySecret == "" {
		return "", 500, errors.New("OSS config not set")
	}

	ossClient, err := oss.NewOSSClient(
		s.ossConfig.Endpoint,
		s.ossConfig.AccessKeyID,
		s.ossConfig.AccessKeySecret,
		s.ossConfig.BucketName,
	)
	if err != nil {
		return "", 500, err
	}

	presignedURL, err := ossClient.GetPresignedURL(video.VideoURL, 3600)
	if err != nil {
		return "", 500, err
	}

	return presignedURL, 0, nil
}

// CheckUserCourse 检查用户是否购买了课程
func (s *VideoService) CheckUserCourse(ctx context.Context, userId, courseId int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.UserCourse{}).Where("user_id = ? AND course_id = ?", userId, courseId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
