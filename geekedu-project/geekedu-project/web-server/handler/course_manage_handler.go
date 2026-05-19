package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"
)

// parseId 将字符串ID转换为int64
func parseId(s string) int64 {
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	CoverUrl    string  `json:"coverUrl" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	CategoryId  int64   `json:"categoryId" binding:"required"`
}

// CreateCourse 创建课程（管理员/讲师）
func CreateCourse(c *gin.Context) {
	userId := c.GetInt64("userId")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params: " + err.Error(),
		})
		return
	}

	resp, err := grpcclient.CourseClient.CreateCourse(context.Background(), &pb.CreateCourseRequest{
		Title:       req.Title,
		Description: req.Description,
		CoverUrl:    req.CoverUrl,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		TeacherId:   userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "create course failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     resp.Code,
		"message":  resp.Message,
		"courseId": resp.CourseId,
	})
}

// AddVideoRequest 添加视频请求
type AddVideoRequest struct {
	Title     string `json:"title" binding:"required"`
	VideoUrl  string `json:"videoUrl" binding:"required"`
	Duration  int32  `json:"duration"`
	IsFree    int32  `json:"isFree"`
	ChapterId int64  `json:"chapterId"`
}

// UploadCourseVideo 添加课程视频（管理员/讲师）
func UploadCourseVideo(c *gin.Context) {
	courseId := c.Param("id")

	var req AddVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UploadCourseVideo bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params: " + err.Error(),
		})
		return
	}

	log.Printf("UploadCourseVideo request: courseId=%s, title=%s, chapterId=%d, videoUrl=%s",
		courseId, req.Title, req.ChapterId, req.VideoUrl)

	// 在web-server端验证chapterId
	if req.ChapterId == 0 {
		log.Printf("UploadCourseVideo error: chapterId is 0")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择视频所属章节",
		})
		return
	}

	resp, err := grpcclient.CourseClient.UploadVideo(context.Background(), &pb.UploadVideoRequest{
		CourseId:  parseId(courseId),
		Title:     req.Title,
		VideoUrl:  req.VideoUrl,
		Duration:  req.Duration,
		IsFree:    req.IsFree,
		ChapterId: req.ChapterId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "add video failed: " + err.Error(),
		})
		return
	}

	// 根据响应码返回适当的HTTP状态
	if resp.Code != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"videoId": resp.VideoId,
	})
}

// DeleteCourse 删除课程（管理员/讲师）
func DeleteCourse(c *gin.Context) {
	courseId := c.Param("id")
	userId := c.GetInt64("userId")

	resp, err := grpcclient.CourseClient.DeleteCourse(context.Background(), &pb.DeleteCourseRequest{
		CourseId:  parseId(courseId),
		TeacherId: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "delete course failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// GetCourseVideos 获取课程视频列表
func GetCourseVideos(c *gin.Context) {
	courseId := c.Param("id")

	resp, err := grpcclient.CourseClient.GetCourseVideos(context.Background(), &pb.GetCourseVideosRequest{
		CourseId: parseId(courseId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "get video list failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"data": gin.H{
			"videos": resp.Videos,
		},
	})
}

// CreateChapterRequest 创建章节请求
type CreateChapterRequest struct {
	Title string `json:"title" binding:"required"`
	Sort  int32  `json:"sort"`
}

// CreateChapter 创建课程章节
func CreateChapter(c *gin.Context) {
	courseId := c.Param("id")

	var req CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params: " + err.Error(),
		})
		return
	}

	resp, err := grpcclient.CourseClient.CreateChapter(context.Background(), &pb.CreateChapterRequest{
		CourseId: parseId(courseId),
		Title:    req.Title,
		Sort:     req.Sort,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "create chapter failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      resp.Code,
		"message":   resp.Message,
		"chapterId": resp.ChapterId,
	})
}
