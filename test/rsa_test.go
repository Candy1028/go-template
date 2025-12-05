package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go-template/pkg/util"
	"os"
	"testing"
)

func TestRsa(t *testing.T) {
	PublicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5E4uJ7RdhHah4nDp86wk
Mvdedy2cIngcjvqK8KcHVEpBeWG4mzeFAfQG6MFWCaFgUwl4oOI2l87tXuhGeRlp
DlbYVLtY8dFbdI802uIk5dPg2PA8sCBEQkSYXe9vihE/IGWhEEgvc2bnHcvWO7t8
cSTyYVlXe6yb1zcj77l4pCj67ZVk4WkZTJXadrPZEGiuGN2w81bvNJKU5+Fs+FCf
AaPUcViN89E4th6bxouQtcu1hGMXg8u0uE6aQYAHrcHmgf6xTB+48Uur0V044Zts
7SISqvmTCelxtd5gNYX3yL1ppujfnGsuu/vYjdTI8Kk7dEJpNLnJk+fySwIMgtYi
gwIDAQAB
-----END PUBLIC KEY-----
`

	//pem
	block, _ := pem.Decode([]byte(PublicKey))
	publtckeyany, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Errorf("Parse PubKeyerr:%s", err)
	}
	// 类型推导
	publtckey := publtckeyany.(*rsa.PublicKey)
	// 使用RSA公钥加密
	rsadata, err := rsa.EncryptPKCS1v15(rand.Reader, publtckey, []byte("123456"))
	if err != nil {
		t.Errorf("EncryptPKCS1v15:%s", err)
	}
	// base64 编码
	basedata := base64.StdEncoding.EncodeToString(rsadata)
	// 获取RSA PRIVATE KEY
	fmt.Println(basedata)
	privatekey, err := os.ReadFile("../config/rsa_private_key.pem")
	if err != nil {
		t.Errorf("Read File:%s", err)
	}
	// RSA解密数据
	data, err := util.RSADecrypt(basedata, privatekey)
	if string(data) != "123456" {
		t.Errorf("Data:%s", err)
	}
	fmt.Println()
}
