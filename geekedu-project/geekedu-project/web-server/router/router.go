package router

import (
	"geekedu/web-server/handler"
	"geekedu/web-server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine) {
	// 静态文件服务（本地上传的文件）
	r.Static("/uploads", "./uploads")
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 公开API
	public := r.Group("/api/v1")
	{
		// 认证相关
		public.POST("/auth/login", handler.Login)
		public.POST("/register", handler.Register)

		// 课程相关（公开）
		public.GET("/courses", handler.GetCourseList)
		public.GET("/courses/:id", handler.GetCourseDetail)
		public.GET("/categories", handler.GetCategoryTree)

		// 视频播放（公开但会检查权限 - 免费视频可直接播放，付费视频需要验证）
		public.GET("/player/:video_id", handler.GetVideoPlayURL)
	}

	// 需要认证的API（所有已登录用户）
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth())
	{
		// 用户相关
		auth.GET("/user/info", handler.GetUserInfo)
		auth.PUT("/user/profile", handler.UpdateProfile)
		auth.PUT("/user/password", handler.ChangePassword)

		// 订单相关（学生）
		auth.POST("/orders", handler.CreateOrder)
		auth.GET("/orders", handler.GetOrderList)
		auth.POST("/orders/:orderNo/pay", handler.PayOrder)
		auth.GET("/orders/check/:courseId", handler.CheckPurchase)
	}

	// 讲师/管理员API
	teacher := r.Group("/api/v1")
	teacher.Use(middleware.JWTAuth(), middleware.RequireRole(middleware.RoleTeacher, middleware.RoleAdmin))
	{
		// 课程管理（讲师端）
		teacher.POST("/courses", handler.CreateCourse)
		teacher.PUT("/courses/:id", handler.UpdateCourse)
		teacher.DELETE("/courses/:id", handler.DeleteCourse)
		teacher.POST("/courses/:id/videos", handler.UploadCourseVideo)
		teacher.GET("/courses/:id/videos", handler.GetCourseVideos)
		// 章节管理
		teacher.POST("/courses/:id/chapters", handler.CreateChapter)
	}

	// 仅管理员API
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(), middleware.RequireRole(middleware.RoleAdmin))
	{
		// 用户管理
		admin.GET("/users", handler.GetUserList)
		admin.PUT("/users/:id/role", handler.UpdateUserRole)

		// 课程审核
		admin.PUT("/courses/:id/status", handler.UpdateCourseStatus)

		// 分类管理
		admin.POST("/categories", handler.CreateCategory)
		admin.PUT("/categories/:id", handler.UpdateCategory)
		admin.DELETE("/categories/:id", handler.DeleteCategory)

		// 订单管理
		admin.GET("/orders", handler.GetAllOrders)
	}

	// 文件上传（需要认证）
	upload := r.Group("/api/v1")
	upload.Use(middleware.JWTAuth())
	{
		upload.POST("/upload", handler.UploadFile)
		upload.POST("/sign-url", handler.SignOSSURL) // OSS URL签名
	}
}
