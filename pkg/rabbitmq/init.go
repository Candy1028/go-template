package rabbitmq

import (
	"fmt"
	"github.com/Candy1028/go-template/log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type rabbitMQConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	VHost    string `json:"vhost"`
	Url      string `json:"url"`
}

func (r *rabbitMQConfig) Get() {
	r.Username = viper.GetString("rabbitmq.username")
	r.Password = viper.GetString("rabbitmq.password")
	r.Host = viper.GetString("rabbitmq.host")
	r.Port = viper.GetInt("rabbitmq.port")
	r.VHost = viper.GetString("rabbitmq.vhost")
	r.Url = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		r.Username, r.Password, r.Host, r.Port, r.VHost,
	)
}

var (
	Conn *amqp.Connection
	mu   sync.RWMutex
)

func InitRabbitMQ() {
	cfg := &rabbitMQConfig{}
	cfg.Get()
	initConnection(cfg)
	go connectionMonitor()

}
func initConnection(cfg *rabbitMQConfig) {
	newConn, err := amqp.Dial(cfg.Url)
	if err != nil {
		log.Logger.Error("建立RabbitMq连接失败")
		return
	}
	Conn = newConn
	log.Logger.Info("成功建立初始RabbitMq连接")
}

func connectionMonitor() {
	retryInterval := time.Second * 3
	maxRetryInterval := time.Minute * 5
	for {
		mu.RLock()
		conn := Conn
		mu.RUnlock()
		if conn != nil {
			connClose := conn.NotifyClose(make(chan *amqp.Error))
			err := <-connClose
			if err != nil {
				log.Logger.Error("RabbitMq连接断开")
			}
		}
		for {
			time.Sleep(retryInterval)
			cfg := &rabbitMQConfig{}
			cfg.Get()
			newConn, err := amqp.DialConfig(cfg.Url, amqp.Config{
				Heartbeat: 10 * time.Second,
			})
			if err != nil {
				log.Logger.Error("RabbitMq重新连接失败")
				retryInterval = time.Duration(float64(retryInterval) * 1.5)
				if retryInterval > maxRetryInterval {
					retryInterval = maxRetryInterval
				}
				continue
			}
			mu.Lock()
			if Conn != nil {
				_ = Conn.Close()
			}
			Conn = newConn
			mu.Unlock()
			log.Logger.Info("成功重新建立RabbitMq连接")
			retryInterval = time.Second * 3
			break

		}

	}
}
