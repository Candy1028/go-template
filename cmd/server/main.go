package main

import (
	"context"
	"github.com/Candy1028/go-template/config"
	"github.com/Candy1028/go-template/internal/middleware/system"
	"github.com/Candy1028/go-template/internal/user"
	"github.com/Candy1028/go-template/log"
	"github.com/Candy1028/go-template/pkg/app"
	"github.com/Candy1028/go-template/pkg/database/mysql"
	"github.com/Candy1028/go-template/pkg/database/redis"
	"github.com/Candy1028/go-template/pkg/util"
	"path/filepath"

	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	log.InitLogger()
	mysql.InitMysql()
	redis.InitRedis()
	// 项目启动入口
	appContext := app.NewAppContext(mysql.DB, redis.DB, redis.Context)
	mode := viper.GetString("gin-mode")
	switch mode {
	case gin.DebugMode:
	case gin.ReleaseMode:
	default:
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	r := gin.Default()
	limMit := system.NewIPRateLimiter(15, 30)
	r.Use(system.CorsSetting(), limMit.Middleware())
	us := r.Group("/api/v1")
	{
		user.NewUserRouter(us, appContext)

	}
	go updateRSAKeyPair(appContext.RDB, appContext.RContext)
	addr := viper.GetString("server.addr")
	if err := r.Run(addr); err != nil {
		log.Logger.Fatal("服务启动失败:", err)
	}
}
func updateRSAKeyPair(rdb *redis2.Client, ctx context.Context) {
	privateKey, publicKey, err := util.GenerateRSAKeyPair()
	if err != nil {
		log.Logger.Errorf("Failed to generate RSA key pair: %v", err)
		return
	}

	// 保存私钥到文件
	privateKeyFilePath := filepath.Join(viper.GetString("rsa.path"), viper.GetString("rsa.name"))
	err = util.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		log.Logger.Errorf("Failed to save private key to file: %v", err)
		return
	}

	// 获取公钥的 PEM 格式字符串
	publicKeyPEM, err := util.GetPublicKeyPEM(publicKey)
	if err != nil {
		log.Logger.Errorf("Failed to get public key PEM: %v", err)
		return
	}

	// 将公钥保存到 Redis
	err = util.SavePublicKeyToRedis(rdb, ctx, publicKeyPEM)
	if err != nil {
		log.Logger.Errorf("Failed to save public key to Redis: %v", err)
		return
	}
	log.Logger.Info("New RSA key pair generated. Public key has been saved to Redis.")
}
