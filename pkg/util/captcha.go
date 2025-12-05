package util

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
)

var store = base64Captcha.DefaultMemStore

// DigitCaptcha 获取数字验证码
func DigitCaptcha() (string, string, error) {
	driver := &base64Captcha.DriverDigit{
		Height:   viper.GetInt("captcha.height"),
		Width:    viper.GetInt("captcha.width"),
		Length:   viper.GetInt("captcha.length"),
		MaxSkew:  viper.GetFloat64("captcha.max_skew"),
		DotCount: viper.GetInt("captcha.dot_count"),
	}
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := captcha.Generate()
	return id, b64s, err
}

// VerifyCaptcha 验证验证码是否正确
func VerifyCaptcha(id, answer string) bool {
	return store.Verify(id, answer, true)
}

// VerifyCode 生成数字验证码
func VerifyCode() (int64, error) {
	randCode, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return 0, errors.New("生成验证码失败")
	}
	code := randCode.Int64() + 100000
	return code, nil
}
