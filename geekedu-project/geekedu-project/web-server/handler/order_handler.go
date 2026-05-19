package handler

import (
	"context"
	"net/http"
	"strconv"

	"geekedu/proto/pb"
	"geekedu/web-server/grpcclient"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	var req struct {
		CourseId int64 `json:"courseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	userId := c.GetInt64("userId")

	resp, err := grpcclient.OrderClient.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		UserId:   userId,
		CourseId: req.CourseId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": gin.H{"orderNo": resp.OrderNo, "amount": resp.Amount}})
}

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) {
	userId := c.GetInt64("userId")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	resp, err := grpcclient.OrderClient.GetOrderList(context.Background(), &pb.GetOrderListRequest{
		UserId:   userId,
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": gin.H{"list": resp.List, "total": resp.Total}})
}

// PayOrder 支付订单
func PayOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	var req struct {
		PayType int `json:"payType"`
	}

	// payType 可选，默认为1
	c.ShouldBindJSON(&req)
	if req.PayType == 0 {
		req.PayType = 1
	}

	resp, err := grpcclient.OrderClient.PayOrder(context.Background(), &pb.PayOrderRequest{
		OrderNo: orderNo,
		PayType: int32(req.PayType),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message})
}

// GetAllOrders 管理员获取所有订单
func GetAllOrders(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	resp, err := grpcclient.OrderClient.GetAllOrders(context.Background(), &pb.GetAllOrdersRequest{
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
		Status:   int32(status),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": resp.Code, "message": resp.Message, "data": gin.H{"list": resp.List, "total": resp.Total, "totalIncome": resp.TotalIncome}})
}

// CheckPurchase 检查用户是否已购买课程
func CheckPurchase(c *gin.Context) {
	userId := c.GetInt64("userId")
	courseId, _ := strconv.ParseInt(c.Param("courseId"), 10, 64)

	resp, err := grpcclient.OrderClient.CheckPurchase(context.Background(), &pb.CheckPurchaseRequest{
		UserId:   userId,
		CourseId: courseId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "service call failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"purchased": resp.Purchased}})
}
