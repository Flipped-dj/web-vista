package grpcclient

import (
	"log"

	"geekedu/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	UserClient   pb.UserServiceClient
	CourseClient pb.CourseServiceClient
	OrderClient  pb.OrderServiceClient
	conn         *grpc.ClientConn
)

// InitGRPCClients 初始化gRPC客户端
func InitGRPCClients(grpcServerAddr string) error {
	// 创建gRPC连接
	var err error
	conn, err = grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// 初始化各个服务客户端
	UserClient = pb.NewUserServiceClient(conn)
	CourseClient = pb.NewCourseServiceClient(conn)
	OrderClient = pb.NewOrderServiceClient(conn)

	log.Printf("gRPC连接成功: %s", grpcServerAddr)
	return nil
}

// Close 关闭gRPC连接
func Close() {
	if conn != nil {
		conn.Close()
	}
}
