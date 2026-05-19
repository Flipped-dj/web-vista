package handler

import (
	"context"
	"net/http"
	"strings"

	"geekedu/common/oss"
	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params: " + err.Error(),
		})
		return
	}

	// 调用logic-server gRPC服务进行注册
	resp, err := grpcclient.UserClient.Register(context.Background(), &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "service call failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"data": gin.H{
			"userId": resp.UserId,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	// 调用logic-server验证并生成JWT token
	resp, err := grpcclient.UserClient.Login(context.Background(), &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "service call failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"data": gin.H{
			"token":    resp.Token,
			"userInfo": resp.UserInfo,
		},
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 从JWT中获取用户ID
	userId := c.GetInt64("userId")

	// 调用logic-server获取用户信息
	resp, err := grpcclient.UserClient.GetUserInfo(context.Background(), &pb.GetUserInfoRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "service call failed",
		})
		return
	}

	// 签名头像URL（如果是OSS链接）
	userInfo := resp.UserInfo
	if userInfo != nil && userInfo.Avatar != "" && strings.Contains(userInfo.Avatar, "aliyuncs.com") {
		client := oss.GetDefaultClient()
		if client != nil {
			if signedURL, err := client.GetPresignedURL(userInfo.Avatar, 3600); err == nil {
				userInfo.Avatar = signedURL
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"data":    userInfo,
	})
}

// UpdateProfile 更新用户资料
func UpdateProfile(c *gin.Context) {
	userId := c.GetInt64("userId")

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	// 调用logic-server更新用户信息
	resp, err := grpcclient.UserClient.UpdateProfile(context.Background(), &pb.UpdateProfileRequest{
		UserId:   userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
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
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userId := c.GetInt64("userId")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	// 调用logic-server修改密码
	resp, err := grpcclient.UserClient.ChangePassword(context.Background(), &pb.ChangePasswordRequest{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "service call failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
	})
}
