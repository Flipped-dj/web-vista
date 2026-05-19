package handler

import (
	"net/http"
	"strconv"

	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"

	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表（管理员）
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	resp, err := grpcclient.UserClient.GetUserList(c, &pb.GetUserListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "get user list failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  resp.Users,
			"total": resp.Total,
		},
	})
}

// UpdateUserRole 更新用户角色（管理员）
func UpdateUserRole(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid user ID",
		})
		return
	}

	var req struct {
		Role int `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	resp, err := grpcclient.UserClient.UpdateUserRole(c, &pb.UpdateUserRoleRequest{
		UserId: userId,
		Role:   int32(req.Role),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "update role failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": resp.Message,
	})
}

// UpdateCourseStatus 更新课程状态（管理员审核）
func UpdateCourseStatus(c *gin.Context) {
	courseId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid course ID",
		})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	resp, err := grpcclient.CourseClient.UpdateCourseStatus(c, &pb.UpdateCourseStatusRequest{
		CourseId: courseId,
		Status:   int32(req.Status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "update course status failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": resp.Message,
	})
}

// CreateCategory 创建分类（管理员）
func CreateCategory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentId int64  `json:"parentId"`
		Sort     int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	resp, err := grpcclient.CourseClient.CreateCategory(c, &pb.CreateCategoryRequest{
		Name:     req.Name,
		ParentId: req.ParentId,
		Sort:     int32(req.Sort),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "create category failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id": resp.CategoryId,
		},
	})
}

// UpdateCategory 更新分类（管理员）
func UpdateCategory(c *gin.Context) {
	categoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid category ID",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	resp, err := grpcclient.CourseClient.UpdateCategory(c, &pb.UpdateCategoryRequest{
		CategoryId: categoryId,
		Name:       req.Name,
		Sort:       int32(req.Sort),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "update category failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": resp.Message,
	})
}

// DeleteCategory 删除分类（管理员）
func DeleteCategory(c *gin.Context) {
	categoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid category ID",
		})
		return
	}

	resp, err := grpcclient.CourseClient.DeleteCategory(c, &pb.DeleteCategoryRequest{
		CategoryId: categoryId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "delete category failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": resp.Message,
	})
}

// UpdateCourse 更新课程（讲师/管理员）
func UpdateCourse(c *gin.Context) {
	courseId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid course ID",
		})
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		CoverUrl    string  `json:"coverUrl"`
		Price       float64 `json:"price"`
		CategoryId  int64   `json:"categoryId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	resp, err := grpcclient.CourseClient.UpdateCourse(c, &pb.UpdateCourseRequest{
		CourseId:    courseId,
		Title:       req.Title,
		Description: req.Description,
		CoverUrl:    req.CoverUrl,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "update course failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": resp.Message,
	})
}
