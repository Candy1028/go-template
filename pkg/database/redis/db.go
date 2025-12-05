package redis

import (
	"context"
	"go-template/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	DB      *redis.Client
	Context = context.Background()
	mu      sync.RWMutex // 添加读写锁保护连接实例
)

type mRedis struct {
	addr string
	db   int
	pwd  string
}

func (m *mRedis) Get() {
	m.addr = viper.GetString("db.redis.addr")
	m.db = viper.GetInt("db.redis.db")
	m.pwd = viper.GetString("db.redis.pwd")
}

func InitRedis() {
	re := &mRedis{}
	re.Get()

	// 首次连接尝试
	initialConnect(re)

	// 独立协程持续监控连接
	go connectionMonitor(re)
}

func initialConnect(re *mRedis) {
	mu.Lock()
	defer mu.Unlock()

	log.Logger.Info("尝试初始化Redis连接...")
	client := redis.NewClient(&redis.Options{
		Addr:     re.addr,
		Password: re.pwd,
		DB:       re.db,
	})

	if _, err := client.Ping(Context).Result(); err != nil {
		log.Logger.Errorf("初始化Redis连接失败: %v", err)
		return
	}

	// 关闭旧连接（如果存在）
	if DB != nil {
		_ = DB.Close()
	}
	DB = client
	//app.RDB = client
	//app.Context = Context
	log.Logger.Info("成功建立初始Redis连接")
}

func connectionMonitor(re *mRedis) {
	retryInterval := time.Second * 10
	for {
		time.Sleep(retryInterval)

		mu.RLock()
		currentClient := DB
		mu.RUnlock()

		// 检查现有连接状态
		if currentClient != nil {
			if _, err := currentClient.Ping(Context).Result(); err == nil {
				continue // 连接正常，继续监控
			}
		}

		// 获取最新配置
		re.Get()
		log.Logger.Info("检测到配置变化或连接断开，尝试重新连接...")

		// 创建新连接
		newClient := redis.NewClient(&redis.Options{
			Addr:     re.addr,
			Password: re.pwd,
			DB:       re.db,
		})

		if _, err := newClient.Ping(Context).Result(); err != nil {
			log.Logger.Errorf("redis 重连尝试失败: %v", err)
			retryInterval = time.Duration(float64(retryInterval) * 1.5) // 退避策略
			if retryInterval > time.Minute {
				retryInterval = time.Minute
			}
			continue
		}

		// 原子替换连接实例
		mu.Lock()
		if DB != nil {
			_ = DB.Close()
		}
		DB = newClient
		//app.RDB = newClient
		//app.Context = Context
		mu.Unlock()
		retryInterval = time.Second * 3 // 重置重试间隔
		log.Logger.Info("Redis连接已成功重建")
	}
}
