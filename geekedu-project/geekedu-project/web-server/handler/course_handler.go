package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"geekedu/common/config"
	"geekedu/common/oss"
	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"

	"github.com/gin-gonic/gin"
)

// signOSSUrl 如果是OSS链接则签名URL
func signOSSUrl(url string) string {
	if url == "" {
		return url
	}
	// 检查是否为我们的OSS存储桶链接
	ossConfig := config.DefaultOSSConfig
	if ossConfig == nil || ossConfig.BucketName == "" {
		return url
	}
	bucketDomain := ossConfig.BucketName + "." + ossConfig.Endpoint
	if !strings.Contains(url, bucketDomain) {
		return url // 不是OSS链接，原样返回
	}

	// 创建OSS客户端并生成签名URL
	ossClient, err := oss.NewOSSClient(
		ossConfig.Endpoint,
		ossConfig.AccessKeyID,
		ossConfig.AccessKeySecret,
		ossConfig.BucketName,
	)
	if err != nil {
		return url
	}

	signedURL, err := ossClient.GetPresignedURL(url, 3600)
	if err != nil {
		return url
	}
	return signedURL
}

// GetCourseList 获取课程列表
func GetCourseList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryIdStr := c.Query("categoryId")
	categoryId, _ := strconv.ParseInt(categoryIdStr, 10, 64)

	resp, err := grpcclient.CourseClient.GetCourseList(context.Background(), &pb.GetCourseListRequest{
		PageNum:    int32(pageNum),
		PageSize:   int32(pageSize),
		CategoryId: categoryId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	// 签名OSS图片URL
	for _, course := range resp.List {
		course.Cover = signOSSUrl(course.Cover)
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": gin.H{"list": resp.List, "total": resp.Total}})
}

// GetCourseDetail 获取课程详情
func GetCourseDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	resp, err := grpcclient.CourseClient.GetCourseDetail(context.Background(), &pb.GetCourseDetailRequest{CourseId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	// 只签名课程封面和讲师头像
	// 视频URL不在这里签名，使用 /player/:video_id 接口获取
	if resp.CourseDetail != nil {
		if resp.CourseDetail.CourseInfo != nil {
			resp.CourseDetail.CourseInfo.Cover = signOSSUrl(resp.CourseDetail.CourseInfo.Cover)
		}
		// 签名讲师头像
		if resp.CourseDetail.Teacher != nil {
			resp.CourseDetail.Teacher.Avatar = signOSSUrl(resp.CourseDetail.Teacher.Avatar)
		}
		// 安全措施：清除视频URL防止泄漏
		// 使用 /player/:video_id 获取签名后的视频URL
		for _, chapter := range resp.CourseDetail.Chapters {
			for _, video := range chapter.Videos {
				video.VideoUrl = "" // 清除视频URL
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": resp.CourseDetail})
}

// GetCategoryTree 获取分类树
func GetCategoryTree(c *gin.Context) {
	resp, err := grpcclient.CourseClient.GetCategoryTree(context.Background(), &pb.GetCategoryTreeRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": resp.List})
}
