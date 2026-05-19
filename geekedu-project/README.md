# GeekEdu 在线教育平台

## 项目概述

GeekEdu 是一个基于**微服务架构**的在线视频学习平台，采用前后端分离设计。

### 核心特性

- 🎓 **课程管理**: 支持课程发布、章节管理、视频上传
- 🔐 **付费内容保护**: 基于 OSS 预签名 URL 的视频播放鉴权
- 👥 **用户系统**: 注册/登录、角色权限、个人中心
- 🛒 **订单系统**: 课程购买、订单管理
- 🎬 **视频播放**: 安全的视频流播放方案

### 核心难点

如何利用云存储 (OSS) 和微服务权限控制，实现"付费内容"的保护与分发：
- 未付费用户**无法获取**视频资源
- 付费用户只能获得**有时效性**的播放权限 (3600秒)
- 直接访问 OSS 原始地址提示 **Access Denied**

## 技术栈

| 层级 | 技术选型 |
|------|----------|
| **前端** | React 18 + Ant Design 5 + Vite + Axios |
| **Web 服务** | Go + Gin + JWT |
| **业务服务** | Go + gRPC |
| **数据层** | MySQL 8.0 + GORM |
| **缓存** | Redis 7 |
| **存储** | 阿里云 OSS (私有 Bucket) |
| **部署** | Docker + Docker Compose |

## 项目结构

```
geekedu-project/
├── web-server/             # Web 接口服务 (Gin + JWT)
│   ├── main.go
│   ├── Dockerfile
│   ├── router/             # 路由配置
│   ├── handler/            # 请求处理器
│   │   ├── user_handler.go
│   │   ├── course_handler.go
│   │   ├── order_handler.go
│   │   ├── video_handler.go    # 视频播放鉴权 (核心)
│   │   └── upload_handler.go
│   ├── middleware/         # 中间件 (JWT认证)
│   └── grpcclient/         # gRPC 客户端
│
├── logic-server/           # 业务逻辑服务 (gRPC)
│   ├── main.go
│   ├── Dockerfile
│   ├── service/            # 业务逻辑实现
│   │   ├── user_service.go
│   │   ├── course_service.go   # 含 GetVideoPlayURL (核心)
│   │   └── order_service.go
│   ├── model/              # 数据模型
│   └── database/           # 数据库连接
│
├── common/                 # 公共组件
│   ├── config/             # 配置管理 (OSS配置)
│   ├── jwt/                # JWT 工具
│   └── oss/                # OSS 客户端 (预签名URL生成)
│
├── proto/                  # Protocol Buffers 定义
│   ├── user.proto
│   ├── course.proto        # 含 GetVideoPlayURL RPC
│   ├── order.proto
│   └── pb/                 # 生成的 Go 代码
│
├── frontend/               # 前端代码 (React)
│   ├── src/
│   │   ├── pages/          # 页面组件
│   │   ├── components/     # 公共组件 (VideoPlayer)
│   │   ├── api/            # API 调用
│   │   └── utils/          # 工具函数
│   ├── Dockerfile
│   └── nginx.conf
│
├── deploy/                 # 部署配置
    ├── docker-compose.yaml
    ├── mysql/
    │   └── init.sql        # 数据库初始化脚本
    └── nginx/

```

## 环境准备

### 1. 阿里云 OSS 配置

1. 创建 Bucket，名称全局唯一 (例如 `geekedu-yourname`)
2. **权限必须设置为"私有" (Private)**
3. 创建 AccessKey (推荐使用 RAM 子账号)
4. 授予 `AliyunOSSFullAccess` 权限

### 2. 配置文件

 `common/config/config.go` 中的 OSS 配置：

```go
var DefaultOSSConfig = &OSSConfig{
    Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
    AccessKeyID:     "your-access-key-id",
    AccessKeySecret: "your-access-key-secret",
    BucketName:      "your-bucket-name",
}
```

## 快速启动

### 一键启动

```bash
cd deploy
docker-compose up -d --build
```

### 查看容器状态

```bash
docker ps
```

### 访问地址

| 服务 | 地址 |
|------|------|
| 前端页面 | http://localhost:3000 |
| API 接口 | http://localhost:8080/api/v1 |

## 功能演示

### 1. 用户注册/登录

访问 http://localhost:3000/register 注册账号，或使用默认管理员账号：
- 用户名: `admin`
- 密码: `admin123`

### 2. 课程浏览

首页展示热门课程，点击课程卡片进入详情页

### 3. 购买课程

1. 登录后进入课程详情页
2. 点击"立即购买"按钮
3. 系统自动扣除余额并建立购买关系

### 4. 视频播放 (核心功能)

**播放鉴权流程**：

1. 用户点击视频播放
2. 前端请求 `GET /api/v1/player/:video_id`
3. 后端校验：
   - 免费视频：直接返回签名 URL
   - 付费视频：检查是否已购买
4. 生成 OSS 预签名 URL (有效期 3600 秒)
5. 前端播放器加载视频

**安全保障**：
- OSS Bucket 设置为私有
- 直接访问原始 URL 返回 403
- 签名 URL 有时效性限制

### 5. 管理后台

管理员登录后访问 http://localhost:3000/admin

功能包括：
- 课程管理 (CRUD、章节、视频)
- 用户管理
- 订单管理
- 分类管理
