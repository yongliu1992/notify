package lib

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var log *logrus.Logger

var once sync.Once
var lock = &sync.Mutex{}
var sysFields map[string]interface{}

/*
 * 单例生成日志实例
 * DebugLevel<InfoLevel<WarnLevel<ErrorLevel
 * 低级别显示比自己高级别的日志
 * 返回实例化对象
 */
func GetLogInstance() *logrus.Entry {
	log =logrus.New()
	log.Out = os.Stdout
	registerLogger := log.WithFields(logrus.Fields{
		"Time": time.Now().Format("2006-01-02 15:04:05"),
	})
	log.SetLevel(logrus.TraceLevel)
	return registerLogger
}