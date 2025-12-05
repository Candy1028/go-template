package mysql

import (
	"context"
	"fmt"
	"github.com/Candy1028/go-template/log"
	gormlogger "github.com/Candy1028/go-template/pkg/database/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var (
	DB *gorm.DB
	mu sync.RWMutex
)

type mysqlConfig struct {
	host     string
	port     string
	username string
	password string
	database string
	charset  string
}

func (m *mysqlConfig) Get() {
	m.host = viper.GetString("db.mysql.host")
	m.port = viper.GetString("db.mysql.port")
	m.username = viper.GetString("db.mysql.username")
	m.password = viper.GetString("db.mysql.pwd")
	m.database = viper.GetString("db.mysql.database")
	m.charset = "utf8mb4"
}

func createDSN(cfg *mysqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.username,
		cfg.password,
		cfg.host,
		cfg.port,
		cfg.database,
		cfg.charset)
}

func InitMysql() {
	cfg := &mysqlConfig{}
	cfg.Get()

	// 首次连接尝试
	initialConnect(cfg)
	// 启动独立监控协程
	go connectionMonitor()
}

func initialConnect(cfg *mysqlConfig) {
	mu.Lock()
	defer mu.Unlock()

	log.Logger.Info("尝试初始化MySQL连接...")
	dsn := createDSN(cfg)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.NewGormLogger(logger.Info),
	})
	if err != nil {
		log.Logger.Errorf("MySQL初始化连接失败: %v", err)
		return
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// 关闭旧连接（如果存在）
	if DB != nil {
		if oldDB, err := DB.DB(); err == nil {
			_ = oldDB.Close()
		}
	}

	DB = db
	//app.DB = db
	log.Logger.Info("成功建立初始MySQL连接")
}

func connectionMonitor() {
	retryInterval := time.Second * 3
	maxRetryInterval := time.Minute * 5

	for {
		time.Sleep(retryInterval)

		mu.RLock()
		currentDB := DB
		mu.RUnlock()

		// 检查连接状态
		if currentDB != nil {
			if sqlDB, err := currentDB.DB(); err == nil {
				if err = sqlDB.PingContext(context.Background()); err == nil {
					continue // 连接正常
				}
			}
		}

		// 获取最新配置
		newCfg := &mysqlConfig{}
		newCfg.Get()

		log.Logger.Info("检测到配置变化或连接断开，尝试重新连接...")

		// 创建新连接
		dsn := createDSN(newCfg)
		newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormlogger.NewGormLogger(logger.Info),
		})

		if err != nil {
			log.Logger.Errorf("MySQL重连失败: %v", err)
			retryInterval = time.Duration(float64(retryInterval) * 1.5)
			if retryInterval > maxRetryInterval {
				retryInterval = maxRetryInterval
			}
			continue
		}

		// 配置新连接池
		if sqlDB, err := newDB.DB(); err == nil {
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}

		// 原子替换连接实例
		mu.Lock()
		if DB != nil {
			if oldDB, err := DB.DB(); err == nil {
				_ = oldDB.Close()
			}
		}
		DB = newDB
		//app.DB = newDB
		mu.Unlock()

		retryInterval = time.Second * 3 // 重置间隔
		log.Logger.Info("MySQL连接已成功重建")
	}
}
