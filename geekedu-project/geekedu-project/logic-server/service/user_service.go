package service

import (
	"context"
	"crypto/md5"
	"fmt"

	"geekedu/common/jwt"
	"geekedu/logic-server/database"
	"geekedu/logic-server/model"
	"geekedu/proto/pb"

	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// ptrToString 安全地将*string转换为string
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 1. 检查用户名是否已存在
	var existUser model.User
	err := database.DB.Where("username = ?", req.Username).First(&existUser).Error
	if err == nil {
		return &pb.RegisterResponse{
			Code:    400,
			Message: "",
		}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return &pb.RegisterResponse{
			Code:    500,
			Message: "database query failed",
		}, err
	}

	// 2. 密码加密（MD5）
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	// 3. 处理昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	// 4. 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: nickname,
		Role:     0, // 默认学生角色
		Status:   1, // 正常状态
	}

	if err := database.DB.Create(user).Error; err != nil {
		return &pb.RegisterResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.RegisterResponse{
		Code:    0,
		Message: "",
		UserId:  user.ID,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. 查询用户
	var user model.User
	err := database.DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.LoginResponse{
				Code:    400,
				Message: "",
			}, nil
		}
		return &pb.LoginResponse{
			Code:    500,
			Message: "database query failed",
		}, err
	}

	// 2. 验证密码
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
	if user.Password != hashedPassword {
		return &pb.LoginResponse{
			Code:    400,
			Message: "",
		}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &pb.LoginResponse{
			Code:    403,
			Message: "",
		}, nil
	}

	// 3. 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return &pb.LoginResponse{
			Code:    500,
			Message: "Token",
		}, err
	}

	return &pb.LoginResponse{
		Code:    0,
		Message: "",
		Token:   token,
		UserInfo: &pb.UserInfo{
			Id:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     int32(user.Role),
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	var user model.User
	err := database.DB.Where("id = ?", req.UserId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetUserInfoResponse{
				Code:    404,
				Message: "user not found",
			}, nil
		}
		return &pb.GetUserInfoResponse{
			Code:    500,
			Message: "database query failed",
		}, err
	}

	return &pb.GetUserInfoResponse{
		Code:    0,
		Message: "success",
		UserInfo: &pb.UserInfo{
			Id:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Role:     int32(user.Role),
		},
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	err := database.DB.Model(&model.User{}).Where("id = ?", req.UserId).Updates(updates).Error
	if err != nil {
		return &pb.UpdateProfileResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.UpdateProfileResponse{
		Code:    0,
		Message: "updated",
	}, nil
}

// GetUserList 获取用户列表（管理员）
func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	var users []model.User
	var total int64

	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询总数
	database.DB.Model(&model.User{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return &pb.GetUserListResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	var pbUsers []*pb.UserInfo
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.UserInfo{
			Id:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Role:      int32(u.Role),
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.GetUserListResponse{
		Code:    0,
		Message: "success",
		Users:   pbUsers,
		Total:   total,
	}, nil
}

// UpdateUserRole 更新用户角色（管理员）
func (s *UserService) UpdateUserRole(ctx context.Context, req *pb.UpdateUserRoleRequest) (*pb.UpdateUserRoleResponse, error) {
	// 验证角色值
	if req.Role < 0 || req.Role > 2 {
		return &pb.UpdateUserRoleResponse{
			Code:    400,
			Message: "invalid role value",
		}, nil
	}

	err := database.DB.Model(&model.User{}).Where("id = ?", req.UserId).Update("role", req.Role).Error
	if err != nil {
		return &pb.UpdateUserRoleResponse{
			Code:    500,
			Message: "",
		}, err
	}

	roleNames := []string{"student", "teacher", "admin"}
	return &pb.UpdateUserRoleResponse{
		Code:    0,
		Message: fmt.Sprintf("user role updated to: %s", roleNames[req.Role]),
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// 查询用户
	var user model.User
	err := database.DB.Where("id = ?", req.UserId).First(&user).Error
	if err != nil {
		return &pb.ChangePasswordResponse{
			Code:    404,
			Message: "user not found",
		}, nil
	}

	// 验证旧密码
	oldHashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.OldPassword)))
	if user.Password != oldHashedPassword {
		return &pb.ChangePasswordResponse{
			Code:    400,
			Message: "incorrect old password",
		}, nil
	}

	// 更新新密码
	newHashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.NewPassword)))
	err = database.DB.Model(&model.User{}).Where("id = ?", req.UserId).Update("password", newHashedPassword).Error
	if err != nil {
		return &pb.ChangePasswordResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.ChangePasswordResponse{
		Code:    0,
		Message: "",
	}, nil
}
