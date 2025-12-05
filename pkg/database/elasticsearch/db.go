package elasticsearch

import (
	"context"
	"crypto/tls"
	"github.com/Candy1028/go-template/log"
	"net/http"

	logs "log"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

var (
	Context = context.Background()
	DB      *elastic.Client
	mu      sync.RWMutex
)

type mEs struct {
	addr        string
	username    string
	pwd         string
	errorLogger MyEsLogger
	infoLogger  MyEsLogger
	traceLogger MyEsLogger
	backoff     elastic.Backoff
}

func (m *mEs) Get() {
	m.addr = viper.GetString("db.elasticsearch.addr")
	m.username = viper.GetString("db.elasticsearch.username")
	m.pwd = viper.GetString("db.elasticsearch.pwd")
	m.errorLogger = MyEsLogger{logs.New(os.Stderr, "[ES-ERROR] ", logs.LstdFlags)}
	m.infoLogger = MyEsLogger{logs.New(os.Stdout, "[ES-INFO] ", logs.LstdFlags)}
	m.traceLogger = MyEsLogger{logs.New(os.Stdout, "[ES-TRACE] ", logs.LstdFlags)}
	m.backoff = elastic.NewExponentialBackoff(10*time.Second, 60*time.Second)
}

type MyEsLogger struct {
	*logs.Logger
}

func InitES() {
	es := &mEs{}
	es.Get()
	// 首次连接尝试
	initialConnect(es)
	//go connectionMonitor(es)
}

func initialConnect(es *mEs) {
	mu.Lock()
	defer mu.Unlock()
	log.Logger.Info("尝试初始化Elasticsearch连接...")
	client, err := elastic.NewClient(
		elastic.SetURL(es.addr),
		elastic.SetBasicAuth(es.username, es.pwd),
		elastic.SetHealthcheck(true),
		elastic.SetSniff(false),
		elastic.SetRetrier(elastic.NewBackoffRetrier(es.backoff)),
		elastic.SetErrorLog(&es.errorLogger),
		elastic.SetInfoLog(&es.infoLogger),
		elastic.SetTraceLog(&es.traceLogger),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //	跳过证书验证
			},
		}),
	)
	if err != nil {
		log.Logger.Errorf("初始化Elasticsearch连接失败: %v", err)
		return

	}

	if _, _, err := client.Ping(es.addr).Do(Context); err != nil {
		log.Logger.Errorf("初始化Elasticsearch连接失败: %v", err)
		return
	}
	DB = client
	log.Logger.Info("成功建立初始Elasticsearch连接")
}
func connectionMonitor(es *mEs) {
	retryInterval := time.Second * 3
	maxRetryInterval := time.Minute * 5

	for {
		time.Sleep(retryInterval)
		mu.RLock()
		currentDB := DB
		mu.RUnlock()

		// 检查连接状态
		if currentDB != nil {
			if _, _, err := currentDB.Ping(es.addr).Do(Context); err == nil {
				if retryInterval != time.Second*3 {
					// 重置间隔
					retryInterval = time.Second * 3
					log.Logger.Info("Elasticsearch连接已成功重建")

				}
				continue
			}
		}

		// 获取最新配置
		newCfg := &mEs{}
		newCfg.Get()
		log.Logger.Info("检测到配置变化或连接断开，尝试重新连接...")
		log.Logger.Errorf("Elasticsearch连接断开")
		retryInterval = time.Duration(float64(retryInterval) * 1.5)
		if retryInterval > maxRetryInterval {
			retryInterval = maxRetryInterval
		}

	}
}
