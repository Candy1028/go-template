package util

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"github.com/spf13/viper"
)

type KodoInfo struct {
	accessKey string
	secretKey string
	Bucket    string
}

func (kodo *KodoInfo) Init() {
	kodo.accessKey = viper.GetString("kodo.access_key")
	kodo.secretKey = viper.GetString("kodo.secret_key")
	kodo.Bucket = viper.GetString("kodo.bucket")
}

// UploadCredentials 获取Kodo上传凭证
func UploadCredentials(fileKey, mimeLimit string, fSizeLimit int64, expiry time.Time) (string, error) {
	info := KodoInfo{}
	info.Init()
	mac := credentials.NewCredentials(info.accessKey, info.secretKey)
	putPolicy, err := uptoken.NewPutPolicy(info.Bucket, expiry)
	if err != nil {
		return "", err
	}
	putPolicy.SetMimeLimit(mimeLimit).SetScope(fmt.Sprintf("%s:%s", info.Bucket, fileKey)).SetInsertOnly(1).SetFsizeLimit(fSizeLimit)
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", err
	}
	return upToken, nil
}
