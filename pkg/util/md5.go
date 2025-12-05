package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
func NewMd5String(s string, sign string) string {
	h := md5.New()
	h.Write([]byte(s + sign))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// ValidateMd5 str前端传入的密码 sign数据库查出的加密秘钥 compareStr数据库中存储的密码
func ValidateMd5(str string, sign string, compareStr string) bool {
	t := NewMd5String(str, sign)
	fmt.Println(t)
	return NewMd5String(str, sign) == compareStr
}

func toLower(s string) string {
	return strings.ToLower(s)
}

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateVerifyCode 生成验证码
func GenerateVerifyCode() int {
	return 100000 + globalRand.Intn(900000)
}

func GenerateRandomUsername() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	username := make([]byte, length)
	for i := range username {
		username[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("用户名%s", string(username))
}

func GenerateSalt(length int) string {
	if length <= 0 {
		length = 16 // 默认长度
	}
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, length)

	rand.Seed(time.Now().UnixNano())
	for i := range salt {
		salt[i] = charset[rand.Intn(len(charset))]
	}

	return string(salt)
}
