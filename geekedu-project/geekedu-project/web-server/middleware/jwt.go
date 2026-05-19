package middleware

import (
	"net/http"
	"strings"

	"geekedu/common/jwt"

	"github.com/gin-gonic/gin"
)

// 角色常量
const (
	RoleStudent = 0 // 学生
	RoleTeacher = 1 // 讲师
	RoleAdmin   = 2 // 管理员
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "authentication token not provided",
			})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// 解析token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid authentication token",
			})
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole 角色验证中间件
func RequireRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "cannot get user role",
			})
			c.Abort()
			return
		}

		userRole := role.(int)
		// 管理员拥有所有权限
		if userRole == RoleAdmin {
			c.Next()
			return
		}

		// 检查用户是否拥有所需角色
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "permission denied",
		})
		c.Abort()
	}
}
