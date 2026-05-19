package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// ServerConfig 服务配置
type ServerConfig struct {
	WebPort   int
	GrpcPort  int
	JWTSecret string
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// GetDatabaseDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*DatabaseConfig, *RedisConfig, *ServerConfig, *OSSConfig) {
	// 首先尝试加载 .env 文件
	loadEnvFile(".env")

	dbConfig := &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 3306),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "geekedu"),
	}

	redisConfig := &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvInt("REDIS_DB", 0),
	}

	serverConfig := &ServerConfig{
		WebPort:   getEnvInt("WEB_PORT", 8080),
		GrpcPort:  getEnvInt("GRPC_PORT", 50051),
		JWTSecret: getEnv("JWT_SECRET", "geekedu-secret-key-2026"),
	}

	ossConfig := &OSSConfig{
		Endpoint:        getEnv("OSS_ENDPOINT", "oss-cn-hangzhou.aliyuncs.com"),
		AccessKeyID:     getEnv("OSS_ACCESS_KEY_ID", ""),
		AccessKeySecret: getEnv("OSS_ACCESS_KEY_SECRET", ""),
		BucketName:      getEnv("OSS_BUCKET_NAME", ""),
	}

	return dbConfig, redisConfig, serverConfig, ossConfig
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt 获取整型环境变量
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// loadEnvFile 从 .env 文件加载环境变量
func loadEnvFile(filename string) {
	// 尝试从当前目录和父目录查找 .env 文件
	paths := []string{
		filename,
		filepath.Join("..", filename),
		filepath.Join("..", "..", filename),
	}

	for _, path := range paths {
		if err := loadEnvFromFile(path); err == nil {
			fmt.Printf("✓ 成功加载配置文件: %s\n", path)
			return
		}
	}
	// 如果找不到 .env 文件，仅使用环境变量或默认值
}

// loadEnvFromFile 从文件加载环境变量
func loadEnvFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 忽略注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// 删除引号（如果有）
		value = strings.Trim(value, "\"'")

		// 只设置环境变量如果还未设置
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// 默认配置（仅用于文档说明）
var (
	DefaultDatabaseConfig = &DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "root",
		DBName:   "geekedu",
	}

	DefaultRedisConfig = &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	DefaultServerConfig = &ServerConfig{
		WebPort:   8080,
		GrpcPort:  50051,
		JWTSecret: "geekedu-secret-key-2026",
	}

	// DefaultOSSConfig 从环境变量加载
	DefaultOSSConfig *OSSConfig

	// 标记是否已初始化
	configInitialized = false
)

func init() {
	// 在init阶段先尝试加载.env文件
	loadEnvFile(".env")

	// 初始化OSS配置，从环境变量读取
	DefaultOSSConfig = &OSSConfig{
		Endpoint:        getEnv("OSS_ENDPOINT", "oss-cn-hangzhou.aliyuncs.com"),
		AccessKeyID:     getEnv("OSS_ACCESS_KEY_ID", ""),
		AccessKeySecret: getEnv("OSS_ACCESS_KEY_SECRET", ""),
		BucketName:      getEnv("OSS_BUCKET_NAME", ""),
	}
	configInitialized = true
}
