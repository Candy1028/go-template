package rabbitmq

import (
	"github.com/Candy1028/go-template/log"
	"time"

	"github.com/streadway/amqp"
)

func Consumer(queueName, queueKey, exchange, kind string, args amqp.Table, handler func(delivery *amqp.Delivery) error) {
	retryInterval := time.Second * 3
	maxRetryInterval := time.Minute * 5
	var ch *amqp.Channel
	var err error
	for {
		if ch != nil {
			_ = ch.Close()
			ch = nil
		}
		retryInterval = time.Duration(float64(retryInterval) * 1.5)
		if retryInterval > maxRetryInterval {
			retryInterval = maxRetryInterval
		}
		ch, err = Conn.Channel()
		if err != nil {
			log.Logger.Errorf("创建Channel失败:%v", err)
			time.Sleep(retryInterval)
			continue
		}
		err = ch.ExchangeDeclare(exchange, kind, false, false, false, false, nil)
		if err != nil {
			log.Logger.Errorf("交换机声明失败:%v", err)
			time.Sleep(retryInterval)
			continue
		}
		_, err := ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			args,
		)
		if err != nil {
			log.Logger.Errorf("声明队列失败:%v", err)
			time.Sleep(retryInterval)
			continue
		}
		err = ch.QueueBind(queueName, queueKey, exchange, false, args)
		if err != nil {
			log.Logger.Errorf("绑定队列失败:%v", err)
			time.Sleep(retryInterval)
			continue
		}
		msgS, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			log.Logger.Errorf("队列监听失败:%v", err)
			time.Sleep(retryInterval)
			continue
		}
		retryInterval = time.Second * 3
		for mas := range msgS {
			if err := handler(&mas); err != nil {
				log.Logger.Errorf("消息处理失败:%v", err)
				_ = mas.Reject(false)
				continue
			}
			_ = mas.Ack(false)
		}
		log.Logger.Error("消息通道关闭，重新建立连接")
	}
}
