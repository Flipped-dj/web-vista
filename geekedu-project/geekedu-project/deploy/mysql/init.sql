-- GeekEdu 数据库初始化脚本
-- 设置字符集
SET NAMES utf8mb4;

SET CHARACTER SET utf8mb4;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS geekedu DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE geekedu;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '密码',
    `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
    `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
    `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
    `role` tinyint(1) DEFAULT 0 COMMENT '角色：0-学员 1-讲师 2-管理员',
    `status` tinyint(1) DEFAULT 1 COMMENT '状态：0-禁用 1-正常',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';

-- 分类表
CREATE TABLE IF NOT EXISTS `category` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `name` varchar(50) NOT NULL COMMENT '分类名称',
    `parent_id` bigint(20) DEFAULT 0 COMMENT '父分类ID',
    `sort` int(11) DEFAULT 0 COMMENT '排序',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '分类表';

-- 课程表
CREATE TABLE IF NOT EXISTS `course` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '课程ID',
    `title` varchar(200) NOT NULL COMMENT '课程标题',
    `cover` varchar(255) DEFAULT NULL COMMENT '课程封面',
    `description` text COMMENT '课程描述',
    `category_id` bigint(20) DEFAULT NULL COMMENT '分类ID',
    `teacher_id` bigint(20) DEFAULT NULL COMMENT '讲师ID',
    `price` decimal(10, 2) DEFAULT 0.00 COMMENT '课程价格',
    `original_price` decimal(10, 2) DEFAULT 0.00 COMMENT '原价',
    `student_count` int(11) DEFAULT 0 COMMENT '学习人数',
    `lesson_count` int(11) DEFAULT 0 COMMENT '课时数',
    `duration` int(11) DEFAULT 0 COMMENT '课程时长（秒）',
    `level` tinyint(1) DEFAULT 1 COMMENT '课程级别：1-入门 2-初级 3-中级 4-高级',
    `is_recommend` tinyint(1) DEFAULT 0 COMMENT '是否推荐：0-否 1-是',
    `status` tinyint(1) DEFAULT 1 COMMENT '状态：0-下架 1-上架',
    `sort` int(11) DEFAULT 0 COMMENT '排序',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_teacher_id` (`teacher_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_recommend` (`is_recommend`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '课程表';

-- 章节表
CREATE TABLE IF NOT EXISTS `chapter` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '章节ID',
    `course_id` bigint(20) NOT NULL COMMENT '课程ID',
    `title` varchar(200) NOT NULL COMMENT '章节标题',
    `sort` int(11) DEFAULT 0 COMMENT '排序',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '章节表';

-- 视频表
CREATE TABLE IF NOT EXISTS `video` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '视频ID',
    `course_id` bigint(20) NOT NULL COMMENT '课程ID',
    `chapter_id` bigint(20) DEFAULT NULL COMMENT '章节ID',
    `title` varchar(200) NOT NULL COMMENT '视频标题',
    `duration` int(11) DEFAULT 0 COMMENT '视频时长（秒）',
    `video_url` varchar(500) DEFAULT NULL COMMENT '视频地址',
    `cover` varchar(255) DEFAULT NULL COMMENT '视频封面',
    `is_free` tinyint(1) DEFAULT 0 COMMENT '是否试看：0-否 1-是',
    `sort` int(11) DEFAULT 0 COMMENT '排序',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_chapter_id` (`chapter_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '视频表';

-- 订单表
CREATE TABLE IF NOT EXISTS `order` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` varchar(50) NOT NULL COMMENT '订单号',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `course_id` bigint(20) NOT NULL COMMENT '课程ID',
    `amount` decimal(10, 2) NOT NULL COMMENT '订单金额',
    `pay_type` tinyint(1) DEFAULT NULL COMMENT '支付方式：1-余额 2-支付宝 3-微信',
    `status` tinyint(1) DEFAULT 0 COMMENT '订单状态：0-待支付 1-已支付 2-已取消 3-已退款',
    `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '订单表';

-- 用户课程关联表
CREATE TABLE IF NOT EXISTS `user_course` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `course_id` bigint(20) NOT NULL COMMENT '课程ID',
    `progress` int(11) DEFAULT 0 COMMENT '学习进度（百分比）',
    `learn_time` int(11) DEFAULT 0 COMMENT '学习时长（秒）',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_course` (`user_id`, `course_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_course_id` (`course_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户课程关联表';

-- 只插入管理员账号
-- 密码为 123456 的 MD5 值: e10adc3949ba59abbe56e057f20f883e
INSERT INTO
    `user` (
        `username`,
        `password`,
        `nickname`,
        `role`,
        `status`
    )
VALUES (
        'admin',
        'e10adc3949ba59abbe56e057f20f883e',
        'Admin',
        2,
        1
    );

-- 插入基础分类数据
INSERT INTO
    `category` (`name`, `parent_id`, `sort`)
VALUES ('前端开发', 0, 1),
    ('后端开发', 0, 2),
    ('移动开发', 0, 3),
    ('人工智能', 0, 4),
    ('全栈工程', 0, 5),
    ('运维部署', 0, 6);