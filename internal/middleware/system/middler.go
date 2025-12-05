package system

import (
	"github.com/Candy1028/go-template/pkg/comment/response"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit // 每秒允许的请求数
	b   int        // 突发容量
}

// NewIPRateLimiter 限流
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	if limiter, exists := i.ips[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户真实 IP
		clientIP := c.Request.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.Request.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		limiter := i.GetLimiter(clientIP)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"code": response.TooManyRequests,
				"msg":  response.GetMessage(429),
			})
			return
		}
		c.Next()
	}
}
