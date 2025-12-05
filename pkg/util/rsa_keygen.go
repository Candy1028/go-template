package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-template/log"
	"os"

	"github.com/redis/go-redis/v9"
)

// GenerateRSAKeyPair 生成 RSA 密钥对
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// SavePrivateKeyToFile 将私钥保存到文件
func SavePrivateKeyToFile(privateKey *rsa.PrivateKey, filePath string) error {
	RSAMu.Lock()
	// 写锁
	defer RSAMu.Unlock()
	PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: PrivateKey,
	}
	// 临时文件
	tmpPath := filePath + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Logger.Error(fmt.Sprintf("PrivateKey File Close Error : %v", err))
		}
	}(file)
	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return err
	}
	return os.Rename(tmpPath, filePath)
}

// GetPublicKeyPEM 获取公钥的 PEM 格式字符串
func GetPublicKeyPEM(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return string(pubPEM), nil
}

// SavePublicKeyToRedis 将公钥保存到 Redis
func SavePublicKeyToRedis(rdb *redis.Client, ctx context.Context, publicKeyPEM string) error {
	return rdb.Set(ctx, "rsa:public_key", publicKeyPEM, 0).Err()
}
