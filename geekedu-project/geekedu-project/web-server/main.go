package main

import (
	"log"
	"os"

	"geekedu/common/config"
	"geekedu/common/oss"
	"geekedu/logic-server/database"
	"geekedu/web-server/grpcclient"
	"geekedu/web-server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载数据库配置
	dbConfig, _, _, _ := config.LoadConfig()

	// 初始化数据库连接（视频播放功能需要）
	if err := database.InitDB(
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化OSS默认客户端
	if err := oss.InitDefaultClient(
		config.DefaultOSSConfig.Endpoint,
		config.DefaultOSSConfig.AccessKeyID,
		config.DefaultOSSConfig.AccessKeySecret,
		config.DefaultOSSConfig.BucketName,
	); err != nil {
		log.Printf("OSS客户端初始化失败: %v (头像等功能可能不可用)", err)
	} else {
		log.Println("OSS客户端初始化成功")
	}

	// 获取gRPC服务器地址
	grpcAddr := os.Getenv("GRPC_SERVER")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	// 初始化gRPC客户端
	if err := grpcclient.InitGRPCClients(grpcAddr); err != nil {
		log.Fatalf("gRPC客户端初始化失败: %v", err)
	}
	defer grpcclient.Close()

	r := gin.Default()

	// 注册路由
	router.SetupRouter(r)

	// 启动服务
	log.Println("Web服务启动在 :8080 端口")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
