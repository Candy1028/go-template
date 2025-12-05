package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// 定义一个密钥，用于签名和验证 JWT

var (
	jwtAccessKey  = []byte(viper.GetString("jwt.access_token_secret"))
	jwtRefreshKey = []byte(viper.GetString("jwt.refresh_token_secret"))
)

// Claims 定义 JWT 的声明结构体
type Claims struct {
	UserID uint
	Email  string
	Role   string
	jwt.RegisteredClaims
}
type Token struct {
	AccessToken  string
	RefreshToken string
}

// GenerateToken 生成 JWT Token
func GenerateToken(userid uint, role, email string) (Token, error) {
	var token Token
	var err error
	token.AccessToken, err = GenerateAccessToken(userid, role, email)
	if err != nil {
		return Token{}, err
	}
	token.RefreshToken, err = GenerateRefreshToken(userid, role, email)
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

// GenerateAccessToken 生成 JWT AccessToken
func GenerateAccessToken(userid uint, role, email string) (string, error) {
	expirationTime := time.Now().Add(viper.GetDuration("jwt.access_token_expiry_time"))
	claims := &Claims{
		UserID: userid,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    viper.GetString("jwt.issuer"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtAccessKey)
	if err != nil {
		return "", errors.New("生成token失败")
	}
	return tokenString, nil
}

// GenerateRefreshToken 生成 JWT RefreshToken
func GenerateRefreshToken(userid uint, role, email string) (string, error) {
	expirationTime := time.Now().Add(viper.GetDuration("jwt.refresh_token_expiry_time"))
	claims := &Claims{
		UserID: userid,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtRefreshKey)
	if err != nil {
		return "", errors.New("生成token失败")
	}
	return tokenString, nil
}

// ValidateAccessToken 验证 JWT AccessToken
func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 解析 token
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名失败: %v", token.Header["alg"])
		}
		return jwtAccessKey, nil
	})

	if err != nil {
		return nil, errors.New("权限不足")
	}

	if !tkn.Valid {
		return nil, errors.New("权限不足")
	}

	return claims, nil
}

// ValidateRefreshToken 验证 JWT RefreshToken
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 解析 token
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名失败: %v", token.Header["alg"])
		}
		return jwtRefreshKey, nil
	})

	if err != nil {
		return nil, errors.New("权限不足")
	}

	if !tkn.Valid {
		return nil, errors.New("权限不足")
	}

	return claims, nil
}
