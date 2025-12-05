package system

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsSetting() gin.HandlerFunc {
	config2 := cors.DefaultConfig()
	config2.AllowHeaders = []string{
		"Authorization",    // 用于 JWT 鉴权
		"Content-Type",     // 常见请求体类型（如 JSON）
		"X-Requested-With", // 标识 AJAX 请求
	} // 允许的请求头
	//t := viper.GetStringSlice("cors.ip")
	config2.AllowAllOrigins = true
	config2.AllowOrigins = []string{}
	config2.AllowCredentials = true

	config2.AddAllowMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}...)
	config2.MaxAge = 86400
	return cors.New(config2)
}
