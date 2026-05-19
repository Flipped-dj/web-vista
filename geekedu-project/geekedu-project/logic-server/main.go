package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"geekedu/common/config"
	"geekedu/logic-server/database"
	"geekedu/logic-server/service"
	"geekedu/proto/pb"

	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	dbConfig, _, _, ossConfig := config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 获取gRPC监听端口
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	healthPort := os.Getenv("HEALTH_PORT")
	if healthPort == "" {
		healthPort = "50052"
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		if err := http.ListenAndServe(":"+healthPort, mux); err != nil {
			log.Printf("健康检查服务启动失败: %v", err)
		}
	}()

	// 创建gRPC服务器
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("无法监听端口 %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterUserServiceServer(grpcServer, service.NewUserService())
	pb.RegisterCourseServiceServer(grpcServer, service.NewCourseService(ossConfig))
	pb.RegisterOrderServiceServer(grpcServer, service.NewOrderService(ossConfig))

	log.Printf("Logic Server 启动成功，监听端口: %s", grpcPort)

	// 启动服务
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC服务启动失败: %v", err)
	}
}
