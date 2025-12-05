package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	RSAMu sync.RWMutex
)

// RSADecryptFormFile 从文件获取密钥解密 RSA 加密的数据
func RSADecryptFormFile(ciphertextBase64 string) ([]byte, error) {
	RSAMu.RLock()
	// 读锁
	defer RSAMu.RUnlock()
	private, err := os.ReadFile(filepath.Join(viper.GetString("rsa.path"), viper.GetString("rsa.name")))
	if err != nil {
		return nil, err
	}
	return RSADecrypt(ciphertextBase64, private)
}

// RSADecrypt 解密 RSA 加密的数据
func RSADecrypt(ciphertextBase64 string, privateKey []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %v", err)
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key is invalid")
	}
	pris, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}
	pri, ok := pris.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext)
}
