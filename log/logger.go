package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"time"
)

var Logger *logrus.Logger

func InitLogger() {
	logdir := viper.GetString("log.dir")
	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		if err := os.Mkdir(logdir, 0755); err != nil {
			logrus.Fatalf("创建日志目录失败: %v", err)
		}
	}
	Logger = logrus.New()
	Logger.SetLevel(logrus.InfoLevel)
	//json格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//控制台输出
	constCout := os.Stdout
	infoWriter := getWriter(path.Join(logdir, "info.log"))
	warnWriter := getWriter(path.Join(logdir, "warn.log"))
	errWriter := getWriter(path.Join(logdir, "error.log"))
	//Hook 分级写入
	Logger.AddHook(&LevelHook{
		Writers: map[logrus.Level]io.Writer{
			logrus.InfoLevel:  io.MultiWriter(constCout, infoWriter),
			logrus.WarnLevel:  io.MultiWriter(constCout, warnWriter),
			logrus.ErrorLevel: io.MultiWriter(constCout, errWriter),
		},
		Formatter: Logger.Formatter,
	})
}

func getWriter(basePath string) io.Writer {
	writer, err := rotatelogs.New(
		basePath+".%Y-%m-%d",
		//rotatelogs.WithLinkName(basePath),
		rotatelogs.WithMaxAge(7*24*time.Hour),     //七天后覆盖
		rotatelogs.WithRotationTime(24*time.Hour), //按天分割日志
	)
	if err != nil {
		logrus.Fatalf("初始化日志文件失败: %v", err)
	}
	return writer
}

// LevelHook 现按级别写入不同文件
type LevelHook struct {
	Writers   map[logrus.Level]io.Writer
	Formatter logrus.Formatter
}

func (hook *LevelHook) Fire(entry *logrus.Entry) error {
	writer, ok := hook.Writers[entry.Level]
	if !ok {
		return nil
	}
	msg, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	return err
}

func (hook *LevelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
