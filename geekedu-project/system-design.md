# GeekEdu 在线教育平台 - 系统设计文档

## 1. 系统架构图

![](./pic/系统架构图.png)

## 2. 数据流向图

### 2.1 用户认证流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant W as Web Server
    participant L as Logic Server
    participant DB as MySQL
    
    U->>W: POST /api/v1/auth/login<br/>{username, password}
    W->>L: gRPC: Login()
    L->>DB: 查询用户信息
    DB-->>L: 用户数据
    L->>L: 验证密码
    L->>L: 生成 JWT Token
    L-->>W: {token, userInfo}
    W-->>U: 返回 Token
```

### 2.2 视频播放鉴权流程（核心功能）

```mermaid
sequenceDiagram
    participant U as 用户
    participant W as Web Server
    participant L as Logic Server
    participant DB as MySQL
    participant OSS as 阿里云OSS
    
    U->>W: GET /api/v1/player/:video_id<br/>[携带JWT Token]
    W->>W: 解析 JWT 获取 userId
    W->>L: gRPC: GetVideoPlayURL(videoId, userId)
    L->>DB: 查询视频信息
    DB-->>L: video {courseId, isFree, videoUrl}
    
    alt 免费视频
        L->>OSS: 生成预签名URL (3600s)
        OSS-->>L: signedUrl
        L-->>W: {playUrl: signedUrl}
    else 付费视频
        L->>DB: 检查是否已购买课程
        DB-->>L: 购买记录
        alt 已购买
            L->>OSS: 生成预签名URL (3600s)
            OSS-->>L: signedUrl
            L-->>W: {playUrl: signedUrl}
        else 未购买
            L-->>W: {code: 403, msg: "请先购买课程"}
        end
    end
    
    W-->>U: 返回播放地址或错误
```

### 2.3 课程购买流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant W as Web Server
    participant L as Logic Server
    participant DB as MySQL
    
    U->>W: POST /api/v1/orders<br/>{courseId}
    W->>L: gRPC: CreateOrder()
    L->>DB: 查询课程价格
    L->>DB: 查询用户余额
    
    alt 余额充足
        L->>DB: 扣除余额
        L->>DB: 创建订单记录
        L->>DB: 创建用户课程关联
        L-->>W: {orderId, success}
    else 余额不足
        L-->>W: {code: 400, msg: "余额不足"}
    end
    
    W-->>U: 返回订单结果
```

## 3. API 接口设计文档

### 3.1 接口总览

![](./pic/接口总览.png)

### 3.2 接口详细说明

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 | 否 |
| POST | `/api/v1/auth/login` | 用户登录 | 否 |
| GET | `/api/v1/user/info` | 获取用户信息 | 是 |
| GET | `/api/v1/courses` | 获取课程列表 | 否 |
| GET | `/api/v1/courses/:id` | 获取课程详情 | 否 |
| POST | `/api/v1/courses` | 创建课程 | 是(讲师) |
| **GET** | **`/api/v1/player/:video_id`** | **获取视频播放地址(核心)** | 可选 |
| POST | `/api/v1/orders` | 创建订单 | 是 |
| GET | `/api/v1/user/courses` | 获取已购课程 | 是 |

### 3.3 核心接口示例

#### 获取视频播放地址

**请求**
```http
GET /api/v1/player/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**成功响应**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "videoId": 1,
        "playUrl": "https://bucket.oss-cn-hangzhou.aliyuncs.com/videos/lesson1.mp4?Expires=1704972800&OSSAccessKeyId=xxx&Signature=xxx",
        "expire": 3600
    }
}
```

**无权限响应**
```json
{
    "code": 403,
    "message": "请先购买课程后观看"
}
```

## 4. 数据库设计

### 4.1 ER 图

![](./pic/ER.png)

### 4.2 核心表结构

```sql
-- 用户表
CREATE TABLE `user` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `username` VARCHAR(50) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `nickname` VARCHAR(50) DEFAULT '',
    `role` TINYINT DEFAULT 1 COMMENT '1学员,2管理员',
    `balance` DECIMAL(10,2) DEFAULT 10000.00,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 课程表
CREATE TABLE `course` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `title` VARCHAR(200) NOT NULL,
    `cover` VARCHAR(500) DEFAULT '',
    `description` TEXT,
    `category_id` BIGINT DEFAULT 0,
    `teacher_id` BIGINT NOT NULL,
    `price` DECIMAL(10,2) DEFAULT 0.00,
    `status` TINYINT DEFAULT 0 COMMENT '0待审核,1已上架,2已下架',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_category` (`category_id`),
    INDEX `idx_teacher` (`teacher_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频表
CREATE TABLE `video` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `course_id` BIGINT NOT NULL,
    `chapter_id` BIGINT DEFAULT NULL,
    `title` VARCHAR(200) NOT NULL,
    `video_url` VARCHAR(500) NOT NULL COMMENT 'OSS存储路径',
    `is_free` TINYINT DEFAULT 0 COMMENT '0付费,1免费',
    INDEX `idx_course` (`course_id`),
    INDEX `idx_chapter` (`chapter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户课程关联表
CREATE TABLE `user_course` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `course_id` BIGINT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_user_course` (`user_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 5. 安全设计

### 5.1 视频保护机制

![](./pic/视频保护机制.png)

### 5.2 鉴权逻辑伪代码

```go
func GetVideoPlayURL(videoId, userId int64) (string, error) {
    // 1. 查询视频信息
    video := GetVideoById(videoId)
    
    // 2. 免费视频直接返回
    if video.IsFree == 1 {
        return generatePresignedURL(video.VideoUrl, 3600), nil
    }
    
    // 3. 检查是否课程作者
    course := GetCourseById(video.CourseId)
    if course.TeacherId == userId {
        return generatePresignedURL(video.VideoUrl, 3600), nil
    }
    
    // 4. 检查是否已购买
    if HasPurchased(userId, video.CourseId) {
        return generatePresignedURL(video.VideoUrl, 3600), nil
    }
    
    // 5. 无权限
    return "", errors.New("请先购买课程")
}
```

## 6. 部署架构

### 6.1 容器编排

![](./pic/容器编排.png)

### 6.2 启动顺序

![](./pic/启动顺序.png)

## 7. oss配置截图

OSS 安全与配置：Bucket 私有

![](./pic/oss1.png)

RAM 用户权限页

![](./pic/oss2.png)

## 8. 技术亮点总结

| 模块 | 技术方案 | 亮点 |
|------|----------|------|
| 视频保护 | OSS 私有 Bucket + 预签名 URL | 防止视频盗链，时效性控制 |
| 服务通信 | gRPC | 高性能、强类型、支持流 |
| 认证方案 | JWT | 无状态、可扩展 |
| 前端架构 | React + Hooks | 组件化、状态管理清晰 |
| 部署方案 | Docker Compose | 一键启动、环境一致 |
