package tools

import (
	"errors"
	"regexp"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("邮箱格式不正确")
	}
	return nil
}

// ValidateMobile 验证手机号格式
func ValidateMobile(mobile string) error {
	pattern := `^1[3-9]\d{9}$`
	matched, err := regexp.MatchString(pattern, mobile)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("手机号格式不正确")
	}
	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	if len(password) > 20 {
		return errors.New("密码长度不能超过20位")
	}
	return nil
}
