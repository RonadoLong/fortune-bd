package log

import (
	"time"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"math"
	"os"
	"fmt"
)
var timeFormat = "02/Jan/2006:15:04:05 -0700"

// Logger is the logrus logger handler
func GinLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}

// config logrus log to local filesystem, with file rotation
//func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
//	baseLogPaht := path.Join(logPath, logFileName)
//	writer, err := rotatelogs.New(
//		baseLogPaht+".%Y%m%d%H%M",
//		rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
//		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
//		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
//	)
//	if err != nil {
//		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
//	}
//	lfHook := lfshook.NewHook(lfshook.WriterMap{
//		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
//		logrus.InfoLevel:  writer,
//		logrus.WarnLevel:  writer,
//		logrus.ErrorLevel: writer,
//		logrus.FatalLevel: writer,
//		logrus.PanicLevel: writer,
//	}, &logrus.JSONFormatter{})
//	logrus.AddHook(lfHook)
//}


// config logrus log to amqp
//func ConfigAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
//	hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
//	log.AddHook(hook)
//}
//
//// config logrus log to es
//func ConfigESLogger(esUrl string, esHOst string, index string) {
//	client, err := elastic.NewClient(elastic.SetURL(esUrl))
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	log.AddHook(esHook)
//}

