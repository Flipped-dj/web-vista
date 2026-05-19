package oss

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSClient struct {
	client     *oss.Client
	bucketName string
}

// NewOSSClient 创建OSS客户端
func NewOSSClient(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSClient, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	return &OSSClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile 上传文件（带重试机制）
func (c *OSSClient) UploadFile(reader io.Reader, filename string) (string, error) {
	bucket, err := c.client.Bucket(c.bucketName)
	if err != nil {
		return "", fmt.Errorf("获取Bucket失败: %v", err)
	}

	// 生成唯一文件名：时间戳 + 原文件名
	timestamp := time.Now().Format("20060102150405")
	ext := path.Ext(filename)
	objectKey := fmt.Sprintf("uploads/%s_%s%s", timestamp, filename[:len(filename)-len(ext)], ext)

	// 读取所有数据到缓冲区，以便重试
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件数据失败: %v", err)
	}

	// 上传文件（带重试，最多3次）
	var uploadErr error
	for i := 0; i < 3; i++ {
		uploadErr = bucket.PutObject(objectKey, bytes.NewReader(data))
		if uploadErr == nil {
			break
		}
		// 重试前等待
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}
	if uploadErr != nil {
		return "", fmt.Errorf("上传文件失败: %v", uploadErr)
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("https://%s.%s/%s", c.bucketName, c.client.Config.Endpoint, objectKey)
	return fileURL, nil
}

// DeleteFile 删除文件
func (c *OSSClient) DeleteFile(objectKey string) error {
	bucket, err := c.client.Bucket(c.bucketName)
	if err != nil {
		return fmt.Errorf("获取Bucket失败: %v", err)
	}

	err = bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// GetPresignedURL 生成预签名URL（用于私有资源的临时访问）
// objectKey: OSS中的文件路径或完整URL
// expireSeconds: 有效期（秒），默认3600秒
func (c *OSSClient) GetPresignedURL(objectKey string, expireSeconds int64) (string, error) {
	bucket, err := c.client.Bucket(c.bucketName)
	if err != nil {
		return "", fmt.Errorf("获取Bucket失败: %v", err)
	}

	// 如果传入的是完整URL，提取objectKey
	prefix := fmt.Sprintf("https://%s.%s/", c.bucketName, c.client.Config.Endpoint)
	if len(objectKey) > len(prefix) && objectKey[:len(prefix)] == prefix {
		objectKey = objectKey[len(prefix):]
	}
	// 也处理 http:// 的情况
	httpPrefix := fmt.Sprintf("http://%s.%s/", c.bucketName, c.client.Config.Endpoint)
	if len(objectKey) > len(httpPrefix) && objectKey[:len(httpPrefix)] == httpPrefix {
		objectKey = objectKey[len(httpPrefix):]
	}

	// 生成预签名URL，有效期默认3600秒（1小时）
	if expireSeconds == 0 {
		expireSeconds = 3600
	}

	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	return signedURL, nil
}

// 默认客户端（全局单例）
var defaultClient *OSSClient

// InitDefaultClient 初始化默认OSS客户端
func InitDefaultClient(endpoint, accessKeyID, accessKeySecret, bucketName string) error {
	client, err := NewOSSClient(endpoint, accessKeyID, accessKeySecret, bucketName)
	if err != nil {
		return err
	}
	defaultClient = client
	return nil
}

// GetDefaultClient 获取默认OSS客户端
func GetDefaultClient() *OSSClient {
	return defaultClient
}
