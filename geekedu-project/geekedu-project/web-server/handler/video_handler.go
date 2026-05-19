package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"geekedu/common/jwt"
	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"

	"github.com/gin-gonic/gin"
)

// GetVideoPlayURL 获取视频播放URL（预签名URL）
// 核心功能：播放鉴权
// 免费视频：无需登录
// 付费视频：需要登录并购买
func GetVideoPlayURL(c *gin.Context) {
	videoIDStr := c.Param("video_id")

	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid video ID",
		})
		return
	}

	// 尝试从Header中获取用户ID（可选认证）
	var userID int64 = 0
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != "" && token != authHeader {
		claims, err := jwt.ParseToken(token)
		if err == nil {
			userID = claims.UserID
		}
	}

	// 调用logic-server gRPC服务获取播放URL
	resp, err := grpcclient.CourseClient.GetVideoPlayURL(context.Background(), &pb.GetVideoPlayURLRequest{
		VideoId: videoID,
		UserId:  userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "service call failed: " + err.Error(),
		})
		return
	}

	if resp.Code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"videoId": videoID,
			"playUrl": resp.PlayUrl,
			"expire":  3600,
		},
	})
}
