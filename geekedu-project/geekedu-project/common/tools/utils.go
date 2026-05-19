package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

// MD5 MD5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	return now.Format("20060102150405") + RandomString(6)
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// FormatDuration 格式化时长（秒转为 HH:MM:SS）
func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return time.Date(0, 0, 0, hours, minutes, secs, 0, time.UTC).Format("15:04:05")
}
