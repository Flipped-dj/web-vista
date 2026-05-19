package service

import (
	"context"
	"fmt"
	"time"

	"geekedu/common/config"
	"geekedu/logic-server/database"
	"geekedu/logic-server/model"
	"geekedu/proto/pb"

	"gorm.io/gorm"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	ossConfig *config.OSSConfig
}

func NewOrderService(ossConfig *config.OSSConfig) *OrderService {
	return &OrderService{
		ossConfig: ossConfig,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var course model.Course
	err := database.DB.Where("id = ?", req.CourseId).First(&course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.CreateOrderResponse{Code: 404, Message: "course not found"}, nil
		}
		return &pb.CreateOrderResponse{Code: 500, Message: "query failed"}, err
	}

	var existingOrder model.Order
	err = database.DB.Where("user_id = ? AND course_id = ? AND status = 1", req.UserId, req.CourseId).First(&existingOrder).Error
	if err == nil {
		return &pb.CreateOrderResponse{Code: 400, Message: "you have already purchased this course"}, nil
	}

	orderNo := fmt.Sprintf("ORD%d%d", time.Now().Unix(), req.UserId)
	order := &model.Order{OrderNo: orderNo, UserID: req.UserId, CourseID: req.CourseId, Amount: course.Price, Status: 0}
	if err := database.DB.Create(order).Error; err != nil {
		return &pb.CreateOrderResponse{Code: 500, Message: ""}, err
	}

	return &pb.CreateOrderResponse{Code: 0, Message: "created", OrderNo: orderNo, Amount: course.Price}, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	var order model.Order
	err := database.DB.Where("order_no = ?", req.OrderNo).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.PayOrderResponse{Code: 404, Message: "order not found"}, nil
		}
		return &pb.PayOrderResponse{Code: 500, Message: "query failed"}, err
	}
	if order.Status == 1 {
		return &pb.PayOrderResponse{Code: 400, Message: "order already paid"}, nil
	}

	now := time.Now()
	err = database.DB.Model(&order).Updates(map[string]interface{}{"status": 1, "pay_type": req.PayType, "pay_time": now}).Error
	if err != nil {
		return &pb.PayOrderResponse{Code: 500, Message: ""}, err
	}

	userCourse := &model.UserCourse{UserID: order.UserID, CourseID: order.CourseID}
	database.DB.Create(userCourse)
	database.DB.Model(&model.Course{}).Where("id = ?", order.CourseID).UpdateColumn("student_count", gorm.Expr("student_count + ?", 1))

	return &pb.PayOrderResponse{Code: 0, Message: ""}, nil
}

// GetOrderList 获取订单列表
func (s *OrderService) GetOrderList(ctx context.Context, req *pb.GetOrderListRequest) (*pb.GetOrderListResponse, error) {
	var orders []model.Order
	var total int64
	query := database.DB.Model(&model.Order{}).Where("user_id = ?", req.UserId)
	if err := query.Count(&total).Error; err != nil {
		return &pb.GetOrderListResponse{Code: 500, Message: "query failed"}, err
	}
	offset := int((req.PageNum - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).Order("created_at DESC").Find(&orders).Error; err != nil {
		return &pb.GetOrderListResponse{Code: 500, Message: "query failed"}, err
	}
	list := make([]*pb.OrderInfo, len(orders))
	for i, order := range orders {
		payTime := ""
		if order.PayTime != nil {
			payTime = order.PayTime.Format("2006-01-02 15:04:05")
		}
		list[i] = &pb.OrderInfo{Id: order.ID, OrderNo: order.OrderNo, UserId: order.UserID, CourseId: order.CourseID, Amount: order.Amount, PayType: int32(order.PayType), Status: int32(order.Status), PayTime: payTime, CreateTime: order.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	return &pb.GetOrderListResponse{Code: 0, Message: "success", List: list, Total: total}, nil
}

// GetAllOrders 管理员获取所有订单
func (s *OrderService) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	var orders []model.Order
	var total int64
	query := database.DB.Model(&model.Order{})

	// 如果指定了状态筛选
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status-1) // status: 1=pending(0), 2=paid(1)
	}

	if err := query.Count(&total).Error; err != nil {
		return &pb.GetAllOrdersResponse{Code: 500, Message: "query failed"}, err
	}

	// 计算已支付订单的总收入
	var totalIncome float64
	database.DB.Model(&model.Order{}).Where("status = 1").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	offset := int((req.PageNum - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).Order("created_at DESC").Find(&orders).Error; err != nil {
		return &pb.GetAllOrdersResponse{Code: 500, Message: "query failed"}, err
	}

	list := make([]*pb.OrderDetailInfo, len(orders))
	for i, order := range orders {
		// 获取用户名
		var user model.User
		username := ""
		if err := database.DB.Where("id = ?", order.UserID).First(&user).Error; err == nil {
			username = user.Username
		}

		// 获取课程名称
		var course model.Course
		courseTitle := ""
		if err := database.DB.Where("id = ?", order.CourseID).First(&course).Error; err == nil {
			courseTitle = course.Title
		}

		payTime := ""
		if order.PayTime != nil {
			payTime = order.PayTime.Format("2006-01-02 15:04:05")
		}

		list[i] = &pb.OrderDetailInfo{
			Id:          order.ID,
			OrderNo:     order.OrderNo,
			UserId:      order.UserID,
			Username:    username,
			CourseId:    order.CourseID,
			CourseTitle: courseTitle,
			Amount:      order.Amount,
			PayType:     int32(order.PayType),
			Status:      int32(order.Status),
			PayTime:     payTime,
			CreateTime:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.GetAllOrdersResponse{Code: 0, Message: "success", List: list, Total: total, TotalIncome: totalIncome}, nil
}

// CheckPurchase 检查用户是否已购买课程
func (s *OrderService) CheckPurchase(ctx context.Context, req *pb.CheckPurchaseRequest) (*pb.CheckPurchaseResponse, error) {
	var count int64
	err := database.DB.Model(&model.UserCourse{}).Where("user_id = ? AND course_id = ?", req.UserId, req.CourseId).Count(&count).Error
	if err != nil {
		return &pb.CheckPurchaseResponse{Purchased: false}, err
	}
	return &pb.CheckPurchaseResponse{Purchased: count > 0}, nil
}
